package messaging

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

type natsMessagingService struct {
	connection            *nats.Conn
	jsonEncodedConnection *nats.EncodedConn
}

func NewNatsMessagingService(cfg *config_api.MessagingAdapter) (MessagingService, error) {
	certs := nats.ClientCert(os.Getenv("SERVICE_TLS_CERT"), os.Getenv("SERVICE_TLS_KEY"))
	rootCA := nats.RootCAs(os.Getenv("SERVICE_TLS_ROOT_CA"))

	log.Printf("Initializing new nats connection with: %s", cfg)
	conn, err := nats.Connect(cfg.GetNats().Url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(1*time.Second),
		certs, rootCA)

	if err != nil {
		return &natsMessagingService{}, err
	}

	err = conn.Flush()
	if err != nil {
		return &natsMessagingService{}, err
	}

	rtt, err := conn.RTT()
	if err != nil {
		return &natsMessagingService{}, err
	}

	log.Printf("NATS server connection initialized with RTT=%s", rtt)

	jsonEncodedConn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return &natsMessagingService{}, err
	}

	return &natsMessagingService{connection: conn,
		jsonEncodedConnection: jsonEncodedConn}, nil
}

func (svc *natsMessagingService) QueueSubscribe(topic string, group string, handler func(msg interface{})) (MessagingQueueSubscription, error) {
	return svc.jsonEncodedConnection.QueueSubscribe(topic, group, handler)
}

func (svc *natsMessagingService) Publish(topic string, msg interface{}) error {
	return svc.jsonEncodedConnection.Publish(topic, msg)
}
