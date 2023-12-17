package kafkastorage

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ServerKafka struct {
	host string
}

func New(kafkaHost string) *ServerKafka {
	if kafkaHost != "" {
		return &ServerKafka{host: kafkaHost}
	}
	return nil
}

func (s *ServerKafka) Send(v []byte, topic string) error {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": s.host})
	if err != nil {
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          v,
	}, nil)

	return err

}
