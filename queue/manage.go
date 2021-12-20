package queue

import (
	"sts/pkg/db"
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
	AutoMigrate()
	for _, consumer := range consumers {
		consumer := consumer
		go func() {
			mq := MQ.ConsumerConnect()
			_ = mq.Consumer(topic, consumer.Queue, consumer.Job, consumer.Sleep, consumer.Retry, consumer.Timeout)
		}()
	}

}

func NewProducer(topic, queue string, message []byte, delay int32) error {
	mq := MQ.ProducerConnect()
	return mq.Producer(topic, queue, message, delay)
}

func AutoMigrate() {
	_ = db.Orm.AutoMigrate(&FailedJobs{})
	ErrJob = make(chan FailedJobs, 1)
	go func() {
		for {
			select {
			case failedJob := <-ErrJob:
				db.Orm.Save(&failedJob)
			}
		}
	}()
}
