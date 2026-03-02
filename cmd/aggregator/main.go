package main

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/ignoxx/toll-calculator/types"
	"google.golang.org/grpc"
)

const (
	listenAddr = ":3000"
	grpcAddr   = ":3001"
)

func main() {
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)

	go func() {
		if err := makeGRPCTransport(grpcAddr, svc); err != nil {
			slog.Error("failed to start gRPC server", "error", err)
			panic(err)
		}
	}()

	if err := makeHTTPTransport(svc); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		slog.Error("failed to listen on address", "addr", listenAddr, "error", err)
		return err
	}

	server := grpc.NewServer()
	types.RegisterDistanceAggregatorServer(server, NewGRPCServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(svc Aggregator) error {
	slog.Info("starting aggregator", "addr", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	return http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuId := r.URL.Query().Get("obuId")
		obuId = strings.TrimSpace(obuId)
		if obuId == "" {
			writeJSON(w, map[string]string{"status": "error", "message": "missing obuId query parameter"})
			return
		}

		obuID, err := strconv.Atoi(obuId)
		if err != nil {
			writeJSON(w, map[string]string{"status": "error", "message": "invalid obuId query parameter"})
			return
		}

		distance, err := svc.CalculateInvoice(obuID)
		if err != nil {
			slog.Error("failed to get invoice", "err", err)
			writeJSON(w, map[string]string{"status": "error", "message": "failed to get invoice"})
			return
		}

		writeJSON(w, map[string]any{"obuId": obuId, "distance": distance})
	}
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
