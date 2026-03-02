package client

import (
	"context"
	"log"

	"github.com/ignoxx/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	types.DistanceAggregatorClient
}

func NewGRPCClient(endpoint string) *GRPCClient {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create gRPC client: " + err.Error())
	}

	c := types.NewDistanceAggregatorClient(conn)
	return &GRPCClient{
		Endpoint:                 endpoint,
		DistanceAggregatorClient: c,
	}
}

func (c *GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggDistanceReq) error {
	_, err := c.AggregateDistance(ctx, aggReq)
	return err
}
