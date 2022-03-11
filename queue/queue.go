package queue

import (
	"fmt"
	"github.com/hr3685930/pkg/goo"
	"github.com/pkg/errors"
)

type JobBase interface {
	Handler() (queueErr *Error)
}

type Queue interface {
	// Connect connect
	Connect() error
	// ProducerConnect Producer
	ProducerConnect() Queue
	// ConsumerConnect consumer connect
	ConsumerConnect() Queue
	// Producer delay
	Producer(topic, queue string, message []byte, delay int32) error
	// Consumer sleep retry
	Consumer(topic, queue string, job JobBase, sleep, retry, timeout int32) error
	// Err report
	Err(failed FailedJobs)

	Close()
}

type Error struct {
	s     string
	stack string
}

func (qe *Error) Error() string {
	return qe.s
}

func Err(err error) *Error {
	e, ok := err.(interface{ goo.Error })
	var stack string
	if ok {
		stack = e.GetStack()
	}else{
		stack = fmt.Sprintf("%+v\n", errors.New(err.Error()))
	}
	return &Error{
		s:     err.Error(),
		stack: stack,
	}
}

func (qe *Error) GetStack() string {
	return qe.stack
}
