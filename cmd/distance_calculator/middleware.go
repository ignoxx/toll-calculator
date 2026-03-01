package main

import (
	"log/slog"
	"time"

	"github.com/ignoxx/toll-calculator/types"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func (l *LogMiddleware) CalculateDistance(data types.ObuData) (float64, error) {
	start := time.Now()
	defer func() {
		slog.Info("calculating distance",
			"obuId", data.ObuID,
			"lat", data.Lat,
			"long", data.Long,
			"took", time.Since(start),
		)
	}()
	return l.next.CalculateDistance(data)
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}
