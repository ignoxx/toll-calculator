package main

import "log"

const (
	topic              = "obu"
	aggregatorEndpoint = "http://127.0.0.1:3000"
)

// consume message from kafka
// calculate distance
// save in DB?
func main() {
	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)

	consumer, err := NewKafkaConsumer(topic)
	if err != nil {
		log.Fatal(err)
	}

	consumer.Start()
}
