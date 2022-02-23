package queue

import (
	"sync"
)

var QueueStore sync.Map

var MQ Queue

var ErrJob chan FailedJobs

func GetQueueDrive(c string) Queue {
	v, ok := QueueStore.Load(c)
	if ok {
		return v.(Queue)
	}
	return nil
}

type Consumer struct {
	Queue   string
	Job     JobBase
	Sleep   int32
	Retry   int32
	Timeout int32
}

type Consumers []*Consumer

func NewConsumer(topic string, consumers Consumers) {
	ErrJob = make(chan FailedJobs, 1)
	for _, consumer := range consumers {
		consumer := consumer
		go func() {
			mq := MQ.ConsumerConnect()
			_ = mq.Consumer(topic, consumer.Queue, consumer.Job, consumer.Sleep, consumer.Retry, consumer.Timeout)
		}()
	}

}

// NewProducer 生产消息
// @Params topic  kafka is topic|rabbitmq is exchange
// @Params key  kafka is partition key|rabbitmq is queue name
// @Params message send body
// @Params delay 延迟
func NewProducer(topic, key string, message []byte, delay int32) error {
	mq := MQ.ProducerConnect()
	return mq.Producer(topic, key, message, delay)
}
