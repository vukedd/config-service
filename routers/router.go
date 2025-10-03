package routers

import (
	"net/http"

	openapimiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/handlers"
	"github.com/vukedd/config-service/metrics"
	"github.com/vukedd/config-service/middleware"
	"github.com/vukedd/config-service/repositories"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

func HandleRequests(router *mux.Router, limiter *rate.Limiter, consulClient *api.Client, tracer trace.Tracer) http.Handler {
	configurationRepository := repositories.NewConfigurationRepository(consulClient, tracer)
	configurationGroupRepository := repositories.NewConfigurationGroupRepository(consulClient, tracer)

	configurationHandler := handlers.NewConfigurationHandler(configurationRepository, tracer)
	configurationGroupHandler := handlers.NewConfigurationGroupHandler(configurationGroupRepository, configurationRepository, tracer)

	// METRICS ENDPOINT
	router.Path("/metrics").Handler(metrics.MetricsHandler())

	// SWAGGER DOCUMENTATION
	// Setup Swagger YAML endpoint
	router.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.yaml")
	})

	// SwaggerUI
	optionsDevelopers := openapimiddleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := openapimiddleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

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
