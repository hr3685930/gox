package local

import (
	"encoding/json"
	"github.com/hr3685930/pkg/queue"
)

type LocalMQ struct {
	Storage chan map[string][]byte
}

func NewLocalMQ() *LocalMQ {
	return &LocalMQ{
		make(chan map[string][]byte, 1),
	}
}

func (l *LocalMQ) Connect() error {
	return nil
}

func (l *LocalMQ) ProducerConnect() queue.Queue {
	return nil
}

func (l *LocalMQ) ConsumerConnect() queue.Queue {
	return nil
}

func (l *LocalMQ) Producer(topic, queue string, message []byte, delay int32) error {
	msg := map[string][]byte{}
	msg[topic] = message
	l.Storage <- msg
	return nil
}

func (l *LocalMQ) Consumer(topic, queue string, job queue.JobBase, sleep, retry, timeout int32) error {
	forever := make(chan bool)
	go func() {
		for m := range l.Storage {
			_ = json.Unmarshal(m[topic], job)
			_ = job.Handler()
		}
	}()

	<-forever
	return nil
}

func (l *LocalMQ) Err(failed queue.FailedJobs) {

}

func (l *LocalMQ) Close() {

}

func (l *LocalMQ) Ping() error {
	return nil
}
