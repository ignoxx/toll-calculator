package main

import (
	"log/slog"

	"github.com/ignoxx/toll-calculator/types"
)

type ObuID = int

type MemoryStore struct {
	db map[ObuID]types.Distance
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		db: make(map[ObuID]types.Distance),
	}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	slog.Info("Inserting distance", "obu_id", distance.ObuID, "distance", distance.Value)
	m.db[distance.ObuID] = distance
	return nil
}
