package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/pkg/errors"
)

type kafkaEvent struct {
	client.Client
	CloudEventID     string
	CloudEventSource string
	CloudEventType   string
	kafka_sarama.Sender
}

var EventClient sarama.Client

func NewKafkaEvent(topic string) (*kafkaEvent, error) {
	if EventClient == nil {
		return nil, errors.New("kafka client is nil")
	}
	sender, err := kafka_sarama.NewSenderFromClient(EventClient, topic)
	if err != nil {
		return nil, err
	}

	c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}
	return &kafkaEvent{Client: c}, nil
}

func (ke *kafkaEvent) SetCloudEventType(eventName string) {
	ke.CloudEventType = eventName
}

func (ke *kafkaEvent) SetCloudEventID(id string) {
	ke.CloudEventID = id
}

func (ke *kafkaEvent) SetCloudEventSource(source string) {
	ke.CloudEventSource = source
}

func (ke *kafkaEvent) Send(ctx context.Context, obj interface{}) error {
	e := cloudevents.NewEvent()
	e.SetID(ke.CloudEventID)
	e.SetType(ke.CloudEventType)
	e.SetSource(ke.CloudEventSource)
	err := e.SetData(cloudevents.ApplicationJSON, obj)
	if err != nil {
		return err
	}
	if result := ke.Client.Send(
		// Set the producer message key
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(e.Type())),
		e,
	); cloudevents.IsUndelivered(result) {
		return errors.Errorf("%+v\n", result)
	}
	return nil
}
