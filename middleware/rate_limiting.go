package middleware

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimit middleware koji ograničava broj zahteva po vremenu
func RateLimit(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Proveri da li je dostignut limit zahteva u datom vremenskom okviru
			if !limiter.Allow() {
				message := map[string]string{
					"message": "rate limit exceeded",
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				err := json.NewEncoder(w).Encode(message)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			// Ako ima dostupnih tokena, pozovi sledeći handler
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitHandler je wrapper funkcija za individualne handler funkcije (deprecated - koristi RateLimit)
func RateLimitHandler(limiter *rate.Limiter, handlerFunc http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := map[string]string{
				"message": "rate limit exceeded", 
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		handlerFunc(w, r)
	})
}