// Configuration Service API
//
// This is a RESTful API service for managing configurations and configuration groups.
// The service provides endpoints for creating, reading, updating, and deleting
// configurations with support for versioning and grouping.
//
// Schemes: http
// Host: localhost:8000
// BasePath: /
// Version: 1.0.0
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/config"
	"github.com/vukedd/config-service/routers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0" // Use available version
	"golang.org/x/time/rate"
)

func main() {
	cfg := config.GetConfig()

	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerAddress)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)
	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	// Finally, set the tracer that can be used for this package.
	tracer := tp.Tracer("config-service")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// 10 requests on initialization,
	// 12 requests per minute (1 request per 5 seconds)
	limiter := rate.NewLimiter(0.2, 10)

	router := mux.NewRouter()

	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.ConsulAddress
	consulClient, _ := api.NewClient(consulConfig)

	srv := http.Server{
		Addr:    cfg.Address,
		Handler: routers.HandleRequests(router, limiter, consulClient, tracer),
	}

	// Starting the server on a new go-routine instead of the main one because the code bellow
	// this block will never be executed since the go-routine will be used by the server which will
	// listen for requests throughout its lifecycle
	go func() {
		fmt.Println("Listening on port" + cfg.Address)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Stopped listening: " + err.Error())
		}
	}()

	// Since we allocated a new go-routine to run the server, the main go-routine is free
	// and in this case it will be used to configure the graceful shutdown mechanism.

	// shutdown - variable that will store the newly created context which listens for terminate and
	// interrupt signals from the OS,
	// stop - variable that stores a function which will stop the shutdown context
	shutdownContext, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	// The stop function is scheduled to be executed right before the function returns.
	// We need this because when the shutdown context is initialized we are starting a new go-routine
	// and if we stop this function and don't free up used resources we are going to have a go-routine leak
	defer stop()

	// shutdown.Done() returns a channel to the main go-routine since there is no receiver on the other
	// side of the pointer, the channel is there, but it is waiting for the signal. When the signal
	// gets broadcasted it will unblock this go-routine and allow the program to continue executing
	<-shutdownContext.Done()

	fmt.Println("Shutdown signal received. Starting graceful shutdown...")

	// Create a new context with a 5-second timeout for the shutdown process.
	// This gives active go-routines time to finish their work.
	timeoutContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Defer cancel to release the resources associated with the timeout context.
	defer cancel()

	if err := srv.Shutdown(timeoutContext); err != nil {
		log.Fatalf("Stopped shutting down: %s", err.Error())
	}
}

func newExporter(address string) (sdktrace.SpanExporter, error) {
	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(address),
		otlptracehttp.WithURLPath("/v1/traces"),
		otlptracehttp.WithInsecure(), // Use for development only
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}
	return exporter, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Use schemaless resource to avoid version conflicts
	r, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceName("config-service"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
