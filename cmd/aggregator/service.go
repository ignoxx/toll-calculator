package main

import "github.com/ignoxx/toll-calculator/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(ObuID) (types.Invoice, error)
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

func (i *InvoiceAggregator) CalculateInvoice(obuId ObuID) (types.Invoice, error) {
	return i.store.GetInvoice(obuId)
}
