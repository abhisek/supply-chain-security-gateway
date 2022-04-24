package dcs

import (
	"log"
	"sync"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

type eventSubscriptionHandler[T any] func(common_models.DomainEvent[T]) error

type eventSubscription[T any] struct {
	name         string
	topic, group string
	handler      eventSubscriptionHandler[T]
}

var dispatcherWg sync.WaitGroup

func registerSubscriber[T any](msgService messaging.MessagingService, subscriber eventSubscription[T]) (messaging.MessagingQueueSubscription, error) {
	log.Printf("Registering disaptcher name:%s topic:%s group:%s", subscriber.name, subscriber.topic, subscriber.group)

	sub, err := msgService.QueueSubscribe(subscriber.topic, subscriber.group, func(msg interface{}) {
		if event, ok := msg.(common_models.DomainEvent[T]); ok {
			subscriber.handler(event)
		} else {
			log.Printf("Error creating a domain event of type T from event msg")
		}
	})

	if err != nil {
		log.Printf("Error registering queue subscriber: %v", err)
	}

	dispatcherWg.Add(1)
	return sub, err
}

func waitForSubscribers() {
	dispatcherWg.Wait()
}
