package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/ignoxx/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	// TODO: read from ENV
	listenAddr = ":3000"
	grpcAddr   = ":3001"
)

func main() {
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go func() {
		if err := makeGRPCTransport(grpcAddr, svc); err != nil {
			slog.Error("failed to start gRPC server", "error", err)
			panic(err)
		}
	}()

	if err := makeHTTPTransport(listenAddr, svc); err != nil {
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

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	var (
		aggMetricHandler = newHTTPMetricsHandler("aggregate")
		invMetricHandler = newHTTPMetricsHandler("invoice")
		aggregateHandler = makeHTTPHandlerFunc(aggMetricHandler.instrument(handleAggregate(svc)))
		invoiceHandler   = makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
	)
	http.HandleFunc("/invoice", invoiceHandler)
	http.HandleFunc("/aggregate", aggregateHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port ", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
