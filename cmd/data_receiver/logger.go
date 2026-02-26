package main

import (
	"log/slog"
	"time"

	"github.com/ignoxx/toll-calculator/types"
)

type LogProducer struct {
	next DataProducer
}

func NewLogProducer(next DataProducer) DataProducer {
	return &LogProducer{
		next: next,
	}
}

func (p *LogProducer) Produce(obuData types.ObuData) error {
	start := time.Now()
	defer func() {
		slog.Info("producing to kafka",
			"obuId", obuData.ObuID,
			"lat", obuData.Lat,
			"long", obuData.Long,
			"took", time.Since(start),
		)
	}()

	return p.next.Produce(obuData)
}

func (p *LogProducer) Close() {
	p.next.Close()
}
