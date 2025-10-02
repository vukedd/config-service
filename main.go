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
	"github.com/vukedd/config-service/routers"
	"golang.org/x/time/rate"
)

func main() {

	// 10 requests on initialization,
	// 12 requests per minute (1 request per 5 seconds)
	limiter := rate.NewLimiter(0.2, 10)

	router := mux.NewRouter()

	consulConfig := api.DefaultConfig()
	consulConfig.Address = "127.0.0.1:8500" // Or from config
	consulClient, _ := api.NewClient(consulConfig)

	srv := http.Server{
		Addr:    ":8000",
		Handler: routers.HandleRequests(router, limiter, consulClient),
	}

	// Starting the server on a new go-routine instead of the main one because the code bellow
	// this block will never be executed since the go-routine will be used by the server which will
	// listen for requests throughout its lifecycle
	go func() {
		fmt.Println("Listening on port :8000")
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
		log.Fatalf("Stopped shutting down: " + err.Error())
	}
}
