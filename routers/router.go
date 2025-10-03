package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/handlers"
	"github.com/vukedd/config-service/metrics"
	"github.com/vukedd/config-service/middleware"
	"github.com/vukedd/config-service/repositories"
	"golang.org/x/time/rate"
)

func HandleRequests(router *mux.Router, limiter *rate.Limiter, consulClient *api.Client) http.Handler {
	configurationRepository := repositories.NewRepository()
	configurationGroupRepository := repositories.NewConfigurationGroupRepository()

	configurationHandler := handlers.NewConfigurationHandler(configurationRepository)
	configurationGroupHandler := handlers.NewConfigurationGroupHandler(configurationGroupRepository, configurationRepository)

	// METRICS ENDPOINT
	router.Path("/metrics").Handler(metrics.MetricsHandler())

	// BASIC OPERATIONS CONFIGURATIONS
	router.Handle("/configurations", middleware.RateLimit(limiter)(http.HandlerFunc(configurationHandler.FindAll))).Methods("GET")
	router.Handle("/configurations/{id}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationHandler.FindById))).Methods("GET")
	router.Handle("/configurations", middleware.IdempotencyMiddleware(consulClient)(middleware.RateLimit(limiter)(http.HandlerFunc(configurationHandler.Create)))).Methods("POST")
	router.Handle("/configurations/{id}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationHandler.Delete))).Methods("DELETE")

	// VERSIONING OPERATIONS CONFIGURATIONS
	router.Handle("/configuration/{name}/{version}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationHandler.DeleteByNameAndVersion))).Methods("DELETE")
	router.Handle("/configuration/{name}/{version}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationHandler.FindByNameAndVersion))).Methods("GET")

	// BASIC OPERATIONS CONFIGURATION GROUP
	router.Handle("/configurationGroups", middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.FindAll))).Methods("GET")
	router.Handle("/configurationGroups/{id}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.FindById))).Methods("GET")
	router.Handle("/configurationGroups/dto/{id}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.FindByIdToDto))).Methods("GET")
	router.Handle("/configurationGroups/{id}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.Delete))).Methods("DELETE")
	router.Handle("/configurationGroups", middleware.IdempotencyMiddleware(consulClient)(middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.Create)))).Methods("POST")
	router.Handle("/configurationGroups/{id}", middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.Update))).Methods("PUT")

	// VERSIONING OPERATIONS CONFIGURATION GROUP
	router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.FindByNameAndVersion).Methods("GET")
	router.HandleFunc("/configurationGroups/dto/{name}/{version}", configurationGroupHandler.FindByNameAndVersionToDto).Methods("GET")
	router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.DeleteByNameAndVersion).Methods("DELETE")

	// Primeni metrics middleware na ceo router
	return middleware.MetricsMiddleware(router)
}
