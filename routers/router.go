package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/handlers"
	"github.com/vukedd/config-service/repositories"
)

func HandleRequests() http.Handler {
	configurationRepository := repositories.NewRepository()
	configurationGroupRepository := repositories.NewConfigurationGroupRepository()

	configurationHandler := handlers.NewConfigurationHandler(configurationRepository)
	configurationGroupHandler := handlers.NewConfigurationGroupHandler(configurationGroupRepository, configurationRepository)

	router := mux.NewRouter()

	// BASIC OPERATIONS CONFIGURATIONS
	router.HandleFunc("/configurations", configurationHandler.FindAll).Methods("GET")
	router.HandleFunc("/configurations/{id}", configurationHandler.FindById).Methods("GET")
	router.HandleFunc("/configurations", configurationHandler.Create).Methods("POST")
	router.HandleFunc("/configurations/{id}", configurationHandler.Delete).Methods("DELETE")

	// VERSIONING OPERATIONS CONFIGURATIONS
	router.HandleFunc("/configuration/{name}/{version}", configurationHandler.DeleteByNameAndVersion).Methods("DELETE")
	router.HandleFunc("/configuration/{name}/{version}", configurationHandler.FindByNameAndVersion).Methods("GET")

	// BASIC OPERATIONS CONFIGURATION GROUP
	router.HandleFunc("/configurationGroups", configurationGroupHandler.FindAll).Methods("GET")
	router.HandleFunc("/configurationGroups/{id}", configurationGroupHandler.FindById).Methods("GET")
	router.HandleFunc("/configurationGroups/dto/{id}", configurationGroupHandler.FindByIdToDto).Methods("GET")
	router.HandleFunc("/configurationGroups/{id}", configurationGroupHandler.Delete).Methods("DELETE")
	router.HandleFunc("/configurationGroups", configurationGroupHandler.Create).Methods("POST")
	router.HandleFunc("/configurationGroups/{id}", configurationGroupHandler.Update).Methods("PUT")

	// VERSIONING OPERATIONS CONFIGURATION GROUP
	router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.FindByNameAndVersion).Methods("GET")
	router.HandleFunc("/configurationGroups/dto/{name}/{version}", configurationGroupHandler.FindByNameAndVersionToDto).Methods("GET")
	router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.DeleteByNameAndVersion).Methods("DELETE")

	return router
}
