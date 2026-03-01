package main

import (
	"log/slog"
	"time"

	"github.com/ignoxx/toll-calculator/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (p *LogMiddleware) AggregateDistance(distance types.Distance) error {
	start := time.Now()
	defer func() {
		slog.Info("aggregating distance",
			"obuId", distance.ObuID,
			"value", distance.Value,
			"timestamp", distance.Timestamp,
			"took", time.Since(start),
		)
	}()

	return p.next.AggregateDistance(distance)
}
