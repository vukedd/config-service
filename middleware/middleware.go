package middleware

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimit is a function that takes the limiter type as well as a function which has http.ResponseWriter
// and *http.Request as params (handlers) as arguments,
func RateLimit(limiter *rate.Limiter, next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	// Returns a new handler that wraps the original one,
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Checks is the limit of requests in the given timeframe reached,
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
		} else {
			// If there are tokens left, the function which was passed as an argument will be called
			next(w, r)
		}
	})
}
