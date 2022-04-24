package adapters

type MessagingHandlerFunc func(data []byte) error

func StartMessagingListener(topic, group string, handler MessagingHandlerFunc) (<-chan bool, error) {
	waiter := make(chan bool)
	return waiter, nil
}
