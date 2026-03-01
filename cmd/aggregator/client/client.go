package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ignoxx/toll-calculator/types"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) AggregateInvoice(distance types.Distance) error {
	b, err := json.Marshal(distance)
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
