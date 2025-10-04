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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

func HandleRequests(router *mux.Router, limiter *rate.Limiter, consulClient *api.Client, tracer trace.Tracer) http.Handler {
    configurationRepository := repositories.NewConfigurationRepository(consulClient, tracer)
    configurationGroupRepository := repositories.NewConfigurationGroupRepository(consulClient, tracer)

    configurationHandler := handlers.NewConfigurationHandler(configurationRepository, tracer)
    configurationGroupHandler := handlers.NewConfigurationGroupHandler(configurationGroupRepository, configurationRepository, tracer)

    // Add tracing middleware to all routes
    router.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)
            defer span.End()
            
            // Add attributes to the span
            span.SetAttributes(
                attribute.String("service.name", "config-service"),
                attribute.String("http.method", r.Method),
                attribute.String("http.url", r.URL.Path),
                attribute.String("http.user_agent", r.UserAgent()),
            )
            
            // Pass the traced context to the next handler
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })

    // METRICS ENDPOINT
    router.Path("/metrics").Handler(metrics.MetricsHandler())

    // TEST TRACE ENDPOINT
    router.HandleFunc("/test-trace", func(w http.ResponseWriter, r *http.Request) {
        _, span := tracer.Start(r.Context(), "test-endpoint")
        defer span.End()
        
        span.SetAttributes(
            attribute.String("test", "true"),
            attribute.String("service.name", "config-service"),
            attribute.String("endpoint", "test-trace"),
        )
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Trace sent to Jaeger!", "service": "config-service"}`))
    }).Methods("GET")

    // SWAGGER DOCUMENTATION
    // Setup Swagger YAML endpoint
    router.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/x-yaml")
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
    router.Handle("/configurationGroups", middleware.IdempotencyMiddleware(consulClient)(middleware.RateLimit(limiter)(http.HandlerFunc(configurationGroupHandler.DeleteByLabel)))).Methods("DELETE")

    // VERSIONING OPERATIONS CONFIGURATION GROUP
    router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.FindByNameAndVersion).Methods("GET")
    router.HandleFunc("/configurationGroups/dto/{name}/{version}", configurationGroupHandler.FindByNameAndVersionToDto).Methods("GET")
    router.HandleFunc("/configurationGroups/{name}/{version}", configurationGroupHandler.DeleteByNameAndVersion).Methods("DELETE")

    // Apply metrics middleware to the entire router
    return middleware.MetricsMiddleware(router)
}