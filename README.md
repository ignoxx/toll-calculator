# toll-calculator

started this project to learn about microservices in Go and to try out gRPC.
it's a simplified version of a real toll calculation system - think of drivers on a highway (autobahn)
where we track their positions and calculate distance to invoice them for tolls.

basically the data flow goes like this:

obu (our test data producer) -> data_receiver (writes to kafka) -> distance_calculator (consumes from kafka, calculates distance) -> aggregator (aggregates into invoices)

services:
- obu - generates random GPS data and sends it via websocket (could be in reality any device sending to our gateway)
- data_receiver - http server that accepts websocket connections and pushes to kafka
- distance_calculator - kafka consumer, calculates distance traveled
- aggregator - http + grpc server, aggregates distances into invoices
- gateway - simple api gateway to the aggregator

how to run:

start kafka:
```
make run-kafka
```

run any service:
```
make obu
make data_receiver
make distance_calculator
make aggregator
make gateway
```

prometheus metrics available at http://localhost:3000/metrics when aggregator is running.
