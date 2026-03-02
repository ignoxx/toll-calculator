package main

import (
	"errors"
	"log/slog"

	"github.com/ignoxx/toll-calculator/types"
)

type ObuID = int

type Storer interface {
	Insert(distance types.Distance) error
	GetInvoice(obuId ObuID) (types.Invoice, error)
}

type MemoryStore struct {
	db map[ObuID][]types.Distance
}

func NewMemoryStore() Storer {
	return &MemoryStore{
		db: make(map[ObuID][]types.Distance),
	}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	slog.Info("Inserting distance", "obu_id", distance.ObuID, "distance", distance.Value)
	m.db[distance.ObuID] = append(m.db[distance.ObuID], distance)
	return nil
}

func (m *MemoryStore) GetInvoice(obuId ObuID) (types.Invoice, error) {
	distance, ok := m.db[obuId]
	if !ok || len(distance) == 0 {
		return types.Invoice{}, errors.New("invoice not found for obuId")
	}

	sum := 0.0
	for _, d := range distance {
		sum += d.Value
	}

	return types.Invoice{
		ObuID:         obuId,
		TotalDistance: sum,
		TotalAmount:   0.04 * sum,
	}, nil
}
