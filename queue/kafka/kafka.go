package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/aaronjan/hunch"
	"github.com/golang-module/carbon"
	"github.com/hr3685930/pkg/queue"
	"github.com/rfyiamcool/go-timewheel"
	"reflect"
	"strings"
	"time"
)

var message chan *msgFuncOpt

type msgFuncOpt struct {
	ConsumerGroupHandler *consumerGroupHandler
	Sess                 sarama.ConsumerGroupSession
	Claim                sarama.ConsumerGroupClaim
	Msg                  *sarama.ConsumerMessage
}

type Kafka struct {
	Cli            sarama.Client
	Brokers        []string
	ConsumerTopics []string
	ProducerTopic  string
	Prefix         string
}

func NewKafka(urls, prefix string) queue.Queue {
	brokers := strings.Split(urls, ",")
	return &Kafka{Prefix: prefix, Brokers: brokers}
}

type consumerGroupHandler struct {
	k         *Kafka
	Sleep     int32
	Retry     int32
	TimeOut   int32
	GroupID   string
	Job       queue.JobBase
	TimeWheel *timewheel.TimeWheel
}

func (k *Kafka) Connect() error {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	//这里的自动提交，是基于被标记过的消息（sess.MarkMessage(msg, “")）
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	client, err := sarama.NewClient(k.Brokers, config)
	if err != nil {
		return err
	}
	k.Cli = client
	return nil
}

func (k *Kafka) ProducerConnect() queue.Queue {
	return &Kafka{Cli: k.Cli, Prefix: k.Prefix, ProducerTopic: k.ProducerTopic, ConsumerTopics: k.ConsumerTopics, Brokers: k.Brokers}
}

func (k *Kafka) ConsumerConnect() queue.Queue {
	return &Kafka{Cli: k.Cli, Prefix: k.Prefix, ProducerTopic: k.ProducerTopic, ConsumerTopics: k.ConsumerTopics, Brokers: k.Brokers}
}

func (k *Kafka) Topic(topic string) {
	k.ProducerTopic = topic
	k.ConsumerTopics = []string{topic}
}

func (k *Kafka) Producer(topic, queueBaseName string, message []byte, delay int32) error {
	p, err := sarama.NewSyncProducerFromClient(k.Cli)
	if err != nil {
		return err
	}
	k.ProducerTopic = topic
	msg := &sarama.ProducerMessage{}
	msg.Topic = k.ProducerTopic
	// 增加key,hash到同一partition保证顺序消费,但发生rebalance时也不能保证顺序性
	// 避免发生rebalance的方法 1.不允许临时增加组下消费者 2.不允许更改partition数
	if queueBaseName != "" {
		msg.Key = sarama.StringEncoder(topic + "_" + queueBaseName)
	}

	var headers []sarama.RecordHeader
	header := sarama.RecordHeader{
		Key:   []byte("delay"),
		Value: queue.Int32ToBytes(delay),
	}
	headers = append(headers, header)
	msg.Headers = headers
	msg.Value = sarama.ByteEncoder(message)
	_, _, err = p.SendMessage(msg)
	if err != nil {
		return err
	}
	return p.Close()
}

func (k *Kafka) Consumer(topic, queueBaseName string, job queue.JobBase, sleep, retry, timeout int32) error {
	groupID := topic + "_" + queueBaseName + "_" + k.Prefix
	group, err := sarama.NewConsumerGroupFromClient(groupID, k.Cli)
	k.ConsumerTopics = []string{topic}
	if err != nil {
		return err
	}
	message = make(chan *msgFuncOpt, 1)
	go func() {
		for {
			select {
			case opt := <-message:
				ConsumerHandler(opt.ConsumerGroupHandler, opt.Sess, opt.Claim, opt.Msg)
			}
		}
	}()
	ctx := context.Background()
	for { //防止rebalance后结束
		topics := k.ConsumerTopics
		handler := &consumerGroupHandler{
			k:       k,
			Job:     job,
			Retry:   retry,
			Sleep:   sleep,
			TimeOut: timeout,
			GroupID: groupID,
		}
		_ = group.Consume(ctx, topics, handler)
	}
}

func (k *Kafka) Err(failed queue.FailedJobs) {
	queue.ErrJob <- failed
}

func (k *Kafka) Close() {
	_ = k.Cli.Close()
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	tw, _ := timewheel.NewTimeWheel(1*time.Second, 360, timewheel.TickSafeMode())
	c.TimeWheel = tw
	c.TimeWheel.Start()
	return nil
}
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	c.TimeWheel.Stop()
	return nil
}

// 消费者会对应一个或者多个partition
// ConsumeClaim 此方法调用次数 = patation数 此方法需要顺序执行
func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		message <- &msgFuncOpt{c, sess, claim, msg}
	}
	return nil
}

func ConsumerHandler(c *consumerGroupHandler, sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, msg *sarama.ConsumerMessage) {
	err := json.Unmarshal(msg.Value, c.Job)
	if err != nil {
		c.k.ExportErr(queue.Err(err), string(msg.Value), c.GroupID)
		sess.MarkMessage(msg, "")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeOut)*time.Second)
	if c.TimeOut == 0 {
		ctx = context.Background()
	}

	headers := make(map[string]interface{}, 1)
	for _, value := range msg.Headers {
		headers[string(value.Key)] = queue.BytesToInt32(value.Value)
	}

	if delay, ok := headers["delay"].(int32); ok && delay > 0 {
		jsonRes := msg.Value
		// interface copy
		msgHandler := reflect.New(reflect.ValueOf(c.Job).Elem().Type()).Interface().(queue.JobBase)
		_ = c.TimeWheel.Add(time.Duration(delay)*time.Second, func() {
			_, err = hunch.Retry(ctx, int(c.Retry)+1, func(ctx context.Context) (interface{}, error) {
				jsonErr := json.Unmarshal(jsonRes, &msgHandler)
				if jsonErr != nil {
					c.k.ExportErr(queue.Err(jsonErr), string(jsonRes), c.GroupID)
					return nil, nil
				}
				handlerErr := msgHandler.Handler()
				if handlerErr != nil {
					c.k.ExportErr(queue.Err(handlerErr), string(jsonRes), c.GroupID)
					c.TimeWheel.Sleep(time.Duration(c.Sleep) * time.Second)
				}
				return nil, handlerErr
			})
			sess.MarkMessage(msg, "")
		})
		cancel()
		return
	}

	_, err = hunch.Retry(ctx, int(c.Retry)+1, func(ctx context.Context) (interface{}, error) {
		handlerErr := c.Job.Handler()
		if handlerErr != nil {
			c.k.ExportErr(queue.Err(handlerErr), string(msg.Value), c.GroupID)
			c.TimeWheel.Sleep(time.Duration(c.Sleep) * time.Second)
		}
		return nil, handlerErr
	})
	sess.MarkMessage(msg, "")
	cancel()
}

func (k *Kafka) ExportErr(err error, msg, groupID string) {
	e := err.(*queue.Error)
	go k.Err(queue.FailedJobs{
		Connection: "kafka",
		Topic:      k.ConsumerTopics[0],
		Queue:      groupID,
		Message:    msg,
		Exception:  err.Error(),
		Stack:      e.GetStack(),
		FiledAt:    carbon.Now(),
	})
}
