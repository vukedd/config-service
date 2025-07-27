package main

import (
	"github.com/vukedd/config-service/routers"
	"log"
	"net/http"
)

func main() {
	router := routers.HandleRequests()

	log.Println("Server starting on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
