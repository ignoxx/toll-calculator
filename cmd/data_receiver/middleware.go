package main

import (
	"log/slog"
	"time"

	"github.com/ignoxx/toll-calculator/types"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) DataProducer {
	return &LogMiddleware{
		next: next,
	}
}

func (p *LogMiddleware) Produce(obuData types.ObuData) error {
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

func (p *LogMiddleware) Close() {
	p.next.Close()
}
