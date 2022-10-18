package messaging

import (
	"fmt"
	"strings"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

type MessagingQueueSubscription interface {
	Unsubscribe() error
}

type MessagingService interface {
	QueueSubscribe(topic string, group string, handler func(msg interface{})) (MessagingQueueSubscription, error)
	Publish(topic string, msg interface{}) error
}

func NewService(adapter *config_api.MessagingAdapter) (MessagingService, error) {
	switch adapter.Type {
	case config_api.MessagingAdapter_NATS:
		return NewNatsMessagingService(adapter)
	case config_api.MessagingAdapter_KAFKA:
		return NewKafkaProtobufMessagingService(strings.Join(adapter.GetKafka().GetBootstrapServers(), ","),
			adapter.GetKafka().SchemaRegistryUrl)
	default:
		return nil, fmt.Errorf("no messaging adapter for: %s", adapter.Type.String())
	}
}
