package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/metrics"
)

// responseWriter wrapper za ResponseWriter klasu
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

// getEndpointPattern izvlači pattern rute iz zahteva
func getEndpointPattern(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route == nil {
		return r.URL.Path
	}
	
	pathTemplate, err := route.GetPathTemplate()
	if err != nil {
		return r.URL.Path
	}
	
	return pathTemplate
}

// isSuccessfulStatusCode proverava da li je status kod uspešan (2xx ili 3xx)
func isSuccessfulStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 400
}

// getStatusClass vraća klasu status koda (2xx, 3xx, 4xx, 5xx)
func getStatusClass(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500 && statusCode < 600:
		return "5xx"
	default:
		return "unknown"
	}
}

// MetricsMiddleware je middleware funkcija koja beleži metrike za svaki HTTP zahtev
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		endpoint := getEndpointPattern(r)
		method := r.Method

		// Obradi wrapper za response writer
		rw := &responseWriter{w, http.StatusOK}

		// Povećaj brojač aktivnih zahteva
		metrics.HttpRequestsInFlight.WithLabelValues(method, endpoint).Inc()
		defer metrics.HttpRequestsInFlight.WithLabelValues(method, endpoint).Dec()

		// Izvršni originalni handler
		next.ServeHTTP(rw, r)

		// Merimo vreme izvršavanja
		duration := time.Since(start).Seconds()
		statusCode := rw.statusCode
		statusCodeStr := strconv.Itoa(statusCode)
		statusClass := getStatusClass(statusCode)

		// Ažuriraj metrike
		metrics.HttpRequestsTotal.WithLabelValues(method, endpoint, statusCodeStr).Inc()
		metrics.HttpResponseDuration.WithLabelValues(method, endpoint).Observe(duration)
		metrics.HttpRequestRate.WithLabelValues(method, endpoint).Inc()

		// Kategorizuj zahteve kao uspešne ili neuspešne
		if isSuccessfulStatusCode(statusCode) {
			metrics.HttpRequestsSuccessful.WithLabelValues(method, endpoint).Inc()
		} else {
			metrics.HttpRequestsFailed.WithLabelValues(method, endpoint, statusClass).Inc()
		}
	})
}