package messaging

import (
	"log"
	"os"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"

	"github.com/nats-io/nats.go"
)

type natsMessagingService struct {
	connection            *nats.Conn
	jsonEncodedConnection *nats.EncodedConn
	config                *common_config.Config
}

func NewNatsMessagingService(config *common_config.Config) (MessagingService, error) {
	certs := nats.ClientCert(os.Getenv("SERVICE_TLS_CERT"), os.Getenv("SERVICE_TLS_KEY"))

	log.Printf("Initializing new nats connection with: %s", config.Global.Messaging.Url)
	conn, err := nats.Connect(config.Global.Messaging.Url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectHandler(natsReconnectHandler()),
		certs)
	if err != nil {
		return &natsMessagingService{}, err
	}

	conn.Flush()
	log.Printf("NATS server connection initialized")

	jsonEncodedConn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return &natsMessagingService{}, err
	}

	return &natsMessagingService{config: config, connection: conn, jsonEncodedConnection: jsonEncodedConn}, nil
}

func (svc *natsMessagingService) QueueSubscribe(topic string, group string, handler func(msg interface{})) (MessagingQueueSubscription, error) {
	return svc.jsonEncodedConnection.QueueSubscribe(topic, group, handler)
}

func (svc *natsMessagingService) Publish(topic string, msg interface{}) error {
	return svc.jsonEncodedConnection.Publish(topic, msg)
}

func natsReconnectHandler() func(*nats.Conn) {
	return func(conn *nats.Conn) {
		log.Printf("Establishing connection with NATS server")
	}
}
