package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ignoxx/toll-calculator/cmd/aggregator/client"
	"github.com/ignoxx/toll-calculator/types"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	calcService CalculatorServicer
	aggClient   client.Client
	isRunning   bool
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, errors.New("failed to create Kafka consumer: " + err.Error())
	}

	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		return nil, errors.New("failed to subscribe to topic: " + err.Error())
	}

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (k *KafkaConsumer) Start() {
	slog.Info("starting Kafka consumer")
	k.isRunning = true
	k.readMessageLoop()
}

func (k *KafkaConsumer) Close() {
	k.isRunning = false
}

func (k *KafkaConsumer) readMessageLoop() {
	defer k.consumer.Close()

	for k.isRunning {
		msg, err := k.consumer.ReadMessage(-1)

		if err != nil {
			slog.Error("failed to read message from Kafka", "err", err)
			continue
		}

		var data types.ObuData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			fmt.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		distance, err := k.calcService.CalculateDistance(data)
		if err != nil {
			slog.Error("failed to calculate distance", "requestID", data.RequestID, "err", err)
			continue
		}

		req := &types.AggDistanceReq{
			Value:     distance,
			Timestamp: time.Now().UnixNano(),
			ObuId:     int32(data.ObuID),
		}

		slog.Info("calculated distance", "requestID", data.RequestID, "obu_id", data.ObuID, "distance", distance)

		if err := k.aggClient.Aggregate(context.TODO(), req); err != nil {
			slog.Error("failed to send aggregate request", "requestID", data.RequestID, "err", err)
			continue
		}
	}

}
