package client

import (
	"context"

	"github.com/ignoxx/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggDistanceReq) error
}
