package main

import (
	"log"

	"github.com/ignoxx/toll-calculator/cmd/aggregator/client"
)

const (
	topic              = "obu"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewClient(aggregatorEndpoint)

	consumer, err := NewKafkaConsumer(topic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	consumer.Start()
}
