package main

import (
	"context"

	"github.com/ignoxx/toll-calculator/types"
)

type GRPCServer struct {
	types.UnimplementedDistanceAggregatorServer
	svc Aggregator
}

func NewGRPCServer(svc Aggregator) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) AggregateDistance(ctx context.Context, req *types.AggDistanceReq) (*types.None, error) {
	distance := types.Distance{
		ObuID:     int(req.ObuId),
		Value:     req.Value,
		Timestamp: req.Timestamp,
	}

	return nil, s.svc.AggregateDistance(distance)
}
