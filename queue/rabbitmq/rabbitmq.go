package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronjan/hunch"
	"github.com/golang-module/carbon"
	"github.com/hr3685930/pkg/queue"
	"github.com/streadway/amqp"
	"time"
)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	QueueName    string
	Exchange     string
	Key          string
	DLKQueueName string
	DLKExchange  string
	DLKKey       string
	MQUrl        string
	Prefix       string
}

func NewRabbitMQ(user, pass, host, port, vhost, prefix string) queue.Queue {
	mqUrl := "amqp://" + user + ":" + pass + "@" + host + ":" + port + vhost
	return &RabbitMQ{MQUrl: mqUrl, Prefix: prefix}
}

func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.MQUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("ampq connect error %s", err))
	}
	channel, err := conn.Channel()
	if err != nil {
		return errors.New(fmt.Sprintf("ampq channel error %s", err))
	}

	r.conn = conn
	r.channel = channel
	return nil
}

func (r *RabbitMQ) ProducerConnect() queue.Queue {
	channel, err := r.conn.Channel()
	if err != nil {
		amqperr := err.(*amqp.Error)
		if amqperr.Code == amqp.ChannelError {
			_ = r.Connect()
			return &RabbitMQ{MQUrl: r.MQUrl, channel: r.channel, conn: r.conn, Prefix: r.Prefix}
		}
		panic(fmt.Sprintf("ampq channel error %s", err))
	}

	r.channel = channel
	return &RabbitMQ{MQUrl: r.MQUrl, channel: channel, conn: r.conn, Prefix: r.Prefix}
}

func (r *RabbitMQ) ConsumerConnect() queue.Queue {
	channel, err := r.conn.Channel()
	if err != nil {
		amqperr := err.(*amqp.Error)
		if amqperr.Code == amqp.ChannelError {
			_ = r.Connect()
			return &RabbitMQ{MQUrl: r.MQUrl, channel: r.channel, conn: r.conn, Prefix: r.Prefix}
		}
		panic(fmt.Sprintf("ampq channel error %s", err))
	}

	r.channel = channel
	return &RabbitMQ{MQUrl: r.MQUrl, channel: channel, conn: r.conn, Prefix: r.Prefix}
}

func (r *RabbitMQ) Close() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func (r *RabbitMQ) SetExchange(exchange string) {
	r.Exchange = exchange
}

func (r *RabbitMQ) SetQueue(queue string) {
	r.QueueName = queue + "_" + r.Prefix
}

func (r *RabbitMQ) SetProducerKey(exchange, queue string) {
	r.Key = exchange + "." + queue + "." + r.Prefix
}

func (r *RabbitMQ) SetConsumerKey(exchange, queue string) {
	r.Key = exchange + "." + queue + ".*"
}

func (r *RabbitMQ) SetDLKExchange() {
	r.DLKExchange = "delay." + r.Exchange
}

func (r *RabbitMQ) SetDLKQueue(sleep int32) {
	r.DLKQueueName = fmt.Sprintf("delay_%d_%s", sleep, r.QueueName)
}

func (r *RabbitMQ) SetDLKKey(sleep int32) {
	r.DLKKey = fmt.Sprintf("delay-%d.%s", sleep, r.Key)
}

func (r *RabbitMQ) Producer(topic, queueBaseName string, message []byte, delay int32) error {
	if r.Exchange == "" {
		if err := r.CreateExchange(topic); err != nil {
			return err
		}
		r.SetExchange(topic)
	}

	r.SetProducerKey(topic, queueBaseName)
	headers := map[string]interface{}{
		"delay": delay,
	}
	if err := r.Publish(r.Exchange, r.Key, message, headers, 2); err != nil {
		return errors.New("publish error")
	}
	return r.channel.Close()
}

func (r *RabbitMQ) DLK(base queue.JobBase, sleep int32, headers map[string]interface{}) error {
	if r.DLKExchange == "" {
		r.SetDLKExchange()
		if err := r.CreateExchange(r.DLKExchange); err != nil {
			return err
		}
	}

	if r.DLKQueueName == "" {
		r.SetDLKQueue(sleep)
		args := amqp.Table{}
		args["x-dead-letter-exchange"] = r.Exchange
		args["x-dead-letter-routing-key"] = r.Key
		args["x-message-ttl"] = 1000 * sleep
		_, err := r.CreateQueue(r.DLKQueueName, args)
		if err != nil {
			return err
		}
	}

	r.SetDLKKey(sleep)
	if err := r.BindKey(r.DLKExchange, r.DLKQueueName, r.DLKKey); err != nil {
		return err
	}

	message, err := json.Marshal(base)
	if err != nil {
		return err
	}
	if err := r.Publish(r.DLKExchange, r.DLKKey, message, headers, 1); err != nil {
		return errors.New("publish error")
	}
	return nil
}

// Consumer  单进程来保证顺序消费
func (r *RabbitMQ) Consumer(topic, queueBaseName string, base queue.JobBase, sleep, retry, timeout int32) error {
	if r.Exchange == "" {
		r.SetExchange(topic)
		if err := r.CreateExchange(topic); err != nil {
			return err
		}
	}

	if r.QueueName == "" {
		r.SetQueue(queueBaseName)
		_, err := r.CreateQueue(r.QueueName, nil)
		if err != nil {
			return err
		}
	}

	r.SetConsumerKey(topic, queueBaseName)
	if err := r.BindKey(r.Exchange, r.QueueName, r.Key); err != nil {
		return err
	}

	if err := r.SetQos(); err != nil {
		return err
	}

	messages, err := r.Consume()
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			jsonErr := json.Unmarshal(d.Body, base)
			if jsonErr != nil {
				r.ExportErr(queue.Err(jsonErr), d)
			}

			// producer delay
			if pDelay, ok := d.Headers["delay"].(int32); ok && pDelay > 0 {
				DlkErr := r.DLK(base, pDelay, nil)
				if DlkErr != nil {
					r.ExportErr(queue.Err(DlkErr), d)
				}
				ackErr := d.Ack(false)
				if ackErr != nil {
					r.ExportErr(queue.Err(ackErr), d)
				}
				continue
			}

			// timeout
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			if timeout == 0 {
				ctx = context.Background()
			}

			// retry and delay
			_, _ = hunch.Retry(ctx, int(retry)+1, func(ctx context.Context) (interface{}, error) {
				handlerErr := base.Handler()
				if handlerErr != nil {
					r.ExportErr(queue.Err(handlerErr), d)
					time.Sleep(time.Second * time.Duration(sleep))
				}
				return nil, handlerErr
			})

			ackErr := d.Ack(false)
			if ackErr != nil {
				r.ExportErr(queue.Err(ackErr), d)
			}
			cancel()
		}
	}()

	<-forever
	return nil
}

func (r *RabbitMQ) Err(failed queue.FailedJobs) {
	queue.ErrJob <- failed
}

func (r *RabbitMQ) CreateExchange(exchange string) error {
	err := r.channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	return err
}

func (r *RabbitMQ) CreateQueue(queueName string, args amqp.Table) (amqp.Queue, error) {
	return r.channel.QueueDeclare(queueName, true, false, false, false, args)
}

func (r *RabbitMQ) BindKey(exchange, queueName, key string) error {
	return r.channel.QueueBind(queueName, key, exchange, false, nil)
}

func (r *RabbitMQ) SetQos() error {
	return r.channel.Qos(1, 0, false)
}

func (r *RabbitMQ) Publish(exchange, key string, message []byte, headers map[string]interface{}, DeliveryMode uint8) error {
	return r.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: DeliveryMode,
			Timestamp:    time.Now(),
			Headers:      headers,
		})
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(r.QueueName, "", false, false, false, false, nil)
}

func (r *RabbitMQ) ExportErr(err error, d amqp.Delivery) {
	e := err.(*queue.Error)
	go r.Err(queue.FailedJobs{
		Connection: "rabbitmq",
		Topic:      r.Exchange,
		Queue:      r.QueueName,
		Message:    string(d.Body),
		Exception:  err.Error(),
		Stack:      e.GetStack(),
		FiledAt:    carbon.Now(),
	})
}
