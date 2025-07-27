package routers

import (
	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/handlers"
	"github.com/vukedd/config-service/repositories"
)

func HandleRequests() *mux.Router {
	configurationRepository := repositories.NewRepository()
	configurationHandler := handlers.NewConfigurationHandler(configurationRepository)

	router := mux.NewRouter()
	router.HandleFunc("/configurations", configurationHandler.FindAll).Methods("GET")
	router.HandleFunc("/configurations/{id}", configurationHandler.FindById).Methods("GET")
	router.HandleFunc("/configurations", configurationHandler.Create).Methods("POST")
	router.HandleFunc("/configurations/{id}", configurationHandler.Delete).Methods("DELETE")
	
	return router
}
