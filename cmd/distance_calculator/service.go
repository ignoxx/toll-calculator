package main

import (
	"math"

	"github.com/ignoxx/toll-calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.ObuData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data types.ObuData) (float64, error) {
	distance := 0.0
	if len(s.prevPoint) == 2 {
		distance = calculateDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}

	s.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculateDistance(lat1, long1, lat2, long2 float64) float64 {
	// Haversine formula
	const R = 6371e3 // Earth radius in meters
	phi1 := lat1 * (math.Pi / 180)
	phi2 := lat2 * (math.Pi / 180)
	deltaPhi := (lat2 - lat1) * (math.Pi / 180)
	deltaLambda := (long2 - long1) * (math.Pi / 180)

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
