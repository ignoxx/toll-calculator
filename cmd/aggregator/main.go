package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ignoxx/toll-calculator/types"
)

const listenAddr = ":3000"

func main() {
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	if err := makeHTTPTransport(svc); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func makeHTTPTransport(svc Aggregator) error {
	slog.Info("starting aggregator", "addr", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	return http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			slog.Error("failed to decode request body", "err", err)
			writeJSON(w, map[string]string{"status": "error", "message": "invalid request body"})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			slog.Error("failed to aggregate distance", "err", err)
			writeJSON(w, map[string]string{"status": "error", "message": "failed to aggregate distance"})
			return
		}

		writeJSON(w, distance)
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
