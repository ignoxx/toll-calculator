package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ignoxx/toll-calculator/types"
)

var wsUpgrader = websocket.Upgrader{}
var listenAddr = ":30000"
var kafkaTopic = "obu"

type DataReceiver struct {
	producer DataProducer
}

func main() {
	slog.Info("starting data receiver", "addr", listenAddr)

	kafkaProducer, err := NewKafkaProducer(kafkaTopic)
	if err != nil {
		slog.Error("failed to create Kafka producer", "err", err)
		panic(err)
	}

	s := DataReceiver{
		producer: NewLogProducer(kafkaProducer),
	}

	defer s.producer.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.wsUpgrade)

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		slog.Error("failed to start server", "err", err)
	}
}

func (s *DataReceiver) wsUpgrade(w http.ResponseWriter, r *http.Request) {
	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to WS", "err", err)
		return
	}

	defer c.Close()

	for {
		var obuData types.ObuData
		if err := c.ReadJSON(&obuData); err != nil {
			slog.Error("failed to read JSON from WS", "err", err)
			break
		}

		// TODO: validate data

		_ = s.producer.Produce(obuData)

	}
}
