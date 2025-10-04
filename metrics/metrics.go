package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Ukupan broj HTTP zahteva
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_requests_total",
			Help: "Ukupan broj HTTP zahteva za configuration service",
		},
		[]string{"method", "endpoint", "status"})

	// Uspešni zahtevi (2xx, 3xx status kodovi)
	HttpRequestsSuccessful = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_requests_successful_total",
			Help: "Broj uspešnih HTTP zahteva (status 2xx, 3xx)",
		},
		[]string{"method", "endpoint"})

	// Neuspešni zahtevi (4xx, 5xx status kodovi)
	HttpRequestsFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_requests_failed_total",
			Help: "Broj neuspešnih HTTP zahteva (status 4xx, 5xx)",
		},
		[]string{"method", "endpoint", "status_class"})

	// Histogram za vreme odgovora po endpoint-u
	HttpResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "config_service_http_response_duration_seconds",
			Help:    "Vreme odgovora HTTP zahteva u sekundama",
			Buckets: prometheus.DefBuckets, // Default bucketi: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"method", "endpoint"})

	// Trenutno aktivni zahtevi
	HttpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "config_service_http_requests_in_flight",
			Help: "Broj trenutno aktivnih HTTP zahteva",
		},
		[]string{"method", "endpoint"})

	// Brzina zahteva po minutu (rate)
	HttpRequestRate = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_request_rate_per_minute",
			Help: "Broj zahteva u minuti za svaki endpoint",
		},
		[]string{"method", "endpoint"})

	// Lista svih metrika koje će biti registrovane
	metricsList = []prometheus.Collector{
		HttpRequestsTotal,
		HttpRequestsSuccessful,
		HttpRequestsFailed,
		HttpResponseDuration,
		HttpRequestsInFlight,
		HttpRequestRate,
	}

	// Prometheus registry za registraciju metrika
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Registruj sve metrike
	prometheusRegistry.MustRegister(metricsList...)
}

// MetricsHandler vraća HTTP handler za Prometheus metrics endpoint
func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}