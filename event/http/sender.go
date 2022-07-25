package http

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	ce "github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/pkg/errors"
)

type httpEvent struct {
	ce.Event
	client.Client
	endpoint string
}

func NewHTTPEvent(endpoint string) (*httpEvent, error) {
	p, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, err
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}

	e := cloudevents.NewEvent()
	return &httpEvent{Event: e, Client: c, endpoint: endpoint}, nil
}

func (he *httpEvent) SetCloudEventType(topic string) {
	he.Event.SetType(topic)
}

func (he *httpEvent) SetCloudEventID(id string) {
	he.Event.SetID(id)
}

func (he *httpEvent) SetCloudEventSource(source string) {
	he.Event.SetSource(source)
}

func (he *httpEvent) Send(ctx context.Context, obj interface{}) error {
	_ = he.Event.SetData(cloudevents.ApplicationJSON, obj)
	ceCTX := cloudevents.ContextWithTarget(ctx, he.endpoint)
	res := he.Client.Send(ceCTX, he.Event)
	if cloudevents.IsUndelivered(res) {
		return errors.Errorf("%+v\n", res)
	} else {
		var httpResult *cehttp.Result
		cloudevents.ResultAs(res, &httpResult)
		return res
	}
}
