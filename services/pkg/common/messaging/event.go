package messaging

import (
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

type DomainEventHandler[T any] func(model common_models.DomainEvent[T], err error) error
type DomainEventMessagingService[T any] interface {
	PublishCreated(T) error
	PublishUpdated(T) error
	PublishDeleted(T) error
	Subscribe(topic string, group string, handler DomainEventHandler[T]) error
}

type domainEventPublisher[T any] struct {
	topic            string
	builder          common_models.DomainEventBuilder[T]
	messagingService MessagingService
}

func NewDomainEventPublisher[T any](topic string,
	builder common_models.DomainEventBuilder[T],
	msgService MessagingService) DomainEventMessagingService[T] {

	return &domainEventPublisher[T]{topic: topic, messagingService: msgService}
}

func (s *domainEventPublisher[T]) PublishCreated(model T) error {
	return s.messagingService.Publish(s.topic, s.builder.Created(model))
}

func (s *domainEventPublisher[T]) PublishUpdated(model T) error {
	return s.messagingService.Publish(s.topic, s.builder.Updated(model))
}

func (s *domainEventPublisher[T]) PublishDeleted(model T) error {
	return s.messagingService.Publish(s.topic, s.builder.Deleted(model))
}

func (s *domainEventPublisher[T]) Subscribe(topic string, group string, handler DomainEventHandler[T]) error {
	_, err := s.messagingService.QueueSubscribe(topic, group, func(msg interface{}) {
		event, err := s.builder.From(msg)
		handler(event, err)
	})

	return err
}
