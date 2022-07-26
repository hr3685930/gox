package event

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	ce "github.com/cloudevents/sdk-go/v2/event"
)

var SendFn = func(ctx context.Context, msg interface{}, endpoint string, event ce.Event) error {
	return nil
}

type rpcEvent struct {
	ce.Event
	endpoint string
}

func NewRpcEvent(endpoint, eventName string) *rpcEvent {
	e := cloudevents.NewEvent()
	return &rpcEvent{Event: e, endpoint: endpoint}
}

func (he *rpcEvent) SetCloudEventType(topic string) {
	he.Event.SetType(topic)
}

func (he *rpcEvent) SetCloudEventID(id string) {
	he.Event.SetID(id)
}

func (he *rpcEvent) SetCloudEventSource(source string) {
	he.Event.SetSource(source)
}

func (he *rpcEvent) Send(ctx context.Context, obj interface{}) error {
	return SendFn(ctx, obj, he.endpoint, he.Event)
}
