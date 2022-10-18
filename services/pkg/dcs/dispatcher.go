package dcs

import (
	"sync"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

type eventSubscriptionHandler[T any] func(common_models.DomainEvent[T]) error

type eventSubscription[T any] struct {
	name         string
	topic, group string
	handler      eventSubscriptionHandler[T]
}

var dispatcherWg sync.WaitGroup

// Register a subscriber to the messaging service and increment
// wait group. Perform generic event to subscriber specific type
// conversion and invoke subscriber business logic
func registerSubscriber[T any](msgService messaging.MessagingService,
	subscriber eventSubscription[T]) (messaging.MessagingQueueSubscription, error) {

	logger.Infof("Registering disaptcher name:%s topic:%s group:%s",
		subscriber.name, subscriber.topic, subscriber.group)

	sub, err := msgService.QueueSubscribe(subscriber.topic, subscriber.group, func(msg interface{}) {
		var event common_models.DomainEvent[T]
		if err := utils.MapStruct(msg, &event); err == nil {
			subscriber.handler(event)
		} else {
			logger.Infof("Error creating a domain event of type T from event msg: %v", err)
		}
	})

	if err != nil {
		logger.Errorf("Error registering queue subscriber: %v", err)
	}

	dispatcherWg.Add(1)
	return sub, err
}

func waitForSubscribers() {
	logger.Infof("Dispatcher waiting for queue subscriptions to close")
	dispatcherWg.Wait()
}
