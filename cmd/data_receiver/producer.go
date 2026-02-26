package main

import (
	"encoding/json"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ignoxx/toll-calculator/types"
)

type DataProducer interface {
	Produce(types.ObuData) error
	Close()
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// go func() {
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
	// 			} else {
	// 				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
	// 			}
	// 		}
	// 	}
	// }()

	return &KafkaProducer{
		producer: p,
	}, nil
}

func (p *KafkaProducer) Produce(obuData types.ObuData) error {
	obuBytes, err := json.Marshal(obuData)
	if err != nil {
		return errors.New("failed to marshal ObuData to JSON: " + err.Error())
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          obuBytes,
	}, nil)

	if err != nil {
		return errors.New("failed to produce message to Kafka: " + err.Error())
	}

	return nil
}

func (p *KafkaProducer) Close() {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
}
