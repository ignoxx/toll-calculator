package main

import "github.com/ignoxx/toll-calculator/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(distance types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}
