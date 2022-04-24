package messaging

type MessagingQueueSubscription interface {
	Unsubscribe() error
}

type MessagingService interface {
	QueueSubscribe(topic string, group string, handler func(msg interface{})) (MessagingQueueSubscription, error)
	Publish(topic string, msg interface{}) error
}
