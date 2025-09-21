package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/handlers"
	"github.com/vukedd/config-service/middleware"
	"github.com/vukedd/config-service/repositories"
	"golang.org/x/time/rate"
)

func HandleRequests() http.Handler {
	configurationRepository := repositories.NewRepository()
	configurationGroupRepository := repositories.NewConfigurationGroupRepository()

	configurationHandler := handlers.NewConfigurationHandler(configurationRepository)
	configurationGroupHandler := handlers.NewConfigurationGroupHandler(configurationGroupRepository, configurationRepository)

	router := mux.NewRouter()

	// 10 requests on initialization,
	// 12 requests per minute (1 request per 5 seconds)
	limiter := rate.NewLimiter(0.2, 10)

	// BASIC OPERATIONS CONFIGURATIONS
	router.Handle("/configurations", middleware.RateLimit(limiter, configurationHandler.FindAll)).Methods("GET")
	router.Handle("/configurations/{id}", middleware.RateLimit(limiter, configurationHandler.FindById)).Methods("GET")
	router.Handle("/configurations", middleware.RateLimit(limiter, configurationHandler.Create)).Methods("POST")
	router.Handle("/configurations/{id}", middleware.RateLimit(limiter, configurationHandler.Delete)).Methods("DELETE")

	// VERSIONING OPERATIONS CONFIGURATIONS
	router.Handle("/configuration/{name}/{version}", middleware.RateLimit(limiter, configurationHandler.DeleteByNameAndVersion)).Methods("DELETE")
	router.Handle("/configuration/{name}/{version}", middleware.RateLimit(limiter, configurationHandler.FindByNameAndVersion)).Methods("GET")

	// BASIC OPERATIONS CONFIGURATION GROUP
	router.Handle("/configurationGroups", middleware.RateLimit(limiter, configurationGroupHandler.FindAll)).Methods("GET")
	router.Handle("/configurationGroups/{id}", middleware.RateLimit(limiter, configurationGroupHandler.FindById)).Methods("GET")
	router.Handle("/configurationGroups/dto/{id}", middleware.RateLimit(limiter, configurationGroupHandler.FindByIdToDto)).Methods("GET")
	router.Handle("/configurationGroups/{id}", middleware.RateLimit(limiter, configurationGroupHandler.Delete)).Methods("DELETE")
	router.Handle("/configurationGroups", middleware.RateLimit(limiter, configurationGroupHandler.Create)).Methods("POST")
	router.Handle("/configurationGroups/{id}", middleware.RateLimit(limiter, configurationGroupHandler.Update)).Methods("PUT")

	// VERSIONING OPERATIONS CONFIGURATION GROUP
	router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.FindByNameAndVersion).Methods("GET")
	router.HandleFunc("/configurationGroups/dto/{name}/{version}", configurationGroupHandler.FindByNameAndVersionToDto).Methods("GET")
	router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.DeleteByNameAndVersion).Methods("DELETE")

	return router
}
