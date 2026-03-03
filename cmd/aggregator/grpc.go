package main

import (
	"context"

	"github.com/ignoxx/toll-calculator/types"
)

type GRPCAggServer struct {
	types.UnimplementedDistanceAggregatorServer
	svc Aggregator
}

func NewGRPCServer(svc Aggregator) *GRPCAggServer {
	return &GRPCAggServer{svc: svc}
}

func (s *GRPCAggServer) AggregateDistance(ctx context.Context, req *types.AggDistanceReq) (*types.None, error) {
	distance := types.Distance{
		ObuID:     int(req.ObuId),
		Value:     req.Value,
		Timestamp: req.Timestamp,
	}

	return &types.None{}, s.svc.AggregateDistance(distance)
}
