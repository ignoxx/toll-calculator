package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ignoxx/toll-calculator/types"
	"github.com/sirupsen/logrus"
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

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (types.Invoice, error) {
	invReq := types.GetInvoiceReq{
		ObuId: int32(id),
	}
	b, err := json.Marshal(&invReq)
	if err != nil {
		return types.Invoice{}, err
	}
	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", id)
	logrus.Infof("requesting get invoice -> %s", endpoint)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return types.Invoice{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.Invoice{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return types.Invoice{}, fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}
	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return types.Invoice{}, err
	}
	defer resp.Body.Close()
	return inv, nil
}
