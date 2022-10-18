package messaging

import (
	"errors"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
)

type kafkaMessagingService struct {
	producer        *kafka.Producer
	serializer      *protobuf.Serializer
	deliveryChannel chan kafka.Event
}

func NewKafkaProtobufMessagingService(bootstrapServers, schemaRegistryUrl string) (MessagingService, error) {
	log.Printf("Kafka msg service init with bootstrap:%s registry:%s",
		bootstrapServers, schemaRegistryUrl)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}

	registryClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(schemaRegistryUrl))
	if err != nil {
		return nil, err
	}

	protobufSerializer, err := protobuf.NewSerializer(registryClient, serde.ValueSerde, protobuf.NewSerializerConfig())
	if err != nil {
		return nil, err
	}

	messageDeliveryNotificationChan := make(chan kafka.Event, 100000)
	messagingService := &kafkaMessagingService{
		producer:        producer,
		serializer:      protobufSerializer,
		deliveryChannel: messageDeliveryNotificationChan,
	}

	go messagingService.deliveryEventHandler()
	return messagingService, nil
}

func (svc *kafkaMessagingService) deliveryEventHandler() {
	log.Printf("Starting Kafka (protobuf) messaging service delivery event handler")
	for msg := range svc.deliveryChannel {
		m, ok := msg.(*kafka.Message)
		if !ok {
			log.Printf("[ERROR] Failed to cast msg to kafka.Message in delivery channel handler")
			continue
		}

		if m.TopicPartition.Error != nil {
			log.Printf("Failed to deliver msg: %v", m.TopicPartition.Error)
		}
	}

	log.Printf("[ERROR] Kafka msg deliver handler QUIT")
}

func (svc *kafkaMessagingService) QueueSubscribe(topic string, group string, handler func(msg interface{})) (MessagingQueueSubscription, error) {
	return nil, errors.New("queue subscription is not supported yet")
}

func (svc *kafkaMessagingService) Publish(topic string, msg interface{}) error {
	payload, err := svc.serializer.Serialize(topic, msg)
	if err != nil {
		return fmt.Errorf("Failed to serialize payload: %v", err)
	}

	return svc.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{},
	}, svc.deliveryChannel)
}
