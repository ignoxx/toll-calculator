package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ignoxx/toll-calculator/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggDistanceReq) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return errors.New("failed to marshal distance data: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return errors.New("failed to create HTTP request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("failed to send HTTP request: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("received non-OK response: " + resp.Status)
	}

	return nil
}
