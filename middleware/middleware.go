package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/models"
	"golang.org/x/time/rate"
)

// RateLimit is a function that takes the limiter type as well as functions which has http.ResponseWriter
// and *http.Request as params (handlers) as arguments,
func RateLimit(limiter *rate.Limiter, next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	// Returns a new handler that wraps the original one,zz
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

func IdempotencyMiddleware(consulClient *api.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idempotencyKey := r.Header.Get("Idempotency-Key")

			// If no key, just proceed to the next handler.
			if idempotencyKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			kv := consulClient.KV()
			keyPath := fmt.Sprintf("idempotency/%s", idempotencyKey)

			// Check if the key already exists.
			pair, _, err := kv.Get(keyPath, nil)
			if err != nil {
				http.Error(w, "Failed to connect to Consul", http.StatusInternalServerError)
				return
			}

			if pair != nil {
				// Key exists, let's see its state.
				var record models.IdempotencyRecord
				if err := json.Unmarshal(pair.Value, &record); err != nil {
					http.Error(w, "Failed to parse stored record", http.StatusInternalServerError)
					return
				}

				// Record is marked as completed, returned the cached response (status and body)
				if record.Status == "completed" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(record.StatusCode)
					w.Write([]byte(record.Body))
					return
				}

				// Record is marked as in progress, rejecting the request to avoid race condition
				if record.Status == "in-progress" {
					http.Error(w, "Request with this idempotency key is already in progress.", http.StatusConflict)
					return
				}
			}

			// Key does not exist. Create a placeholder (lock).
			placeholder := &models.IdempotencyRecord{Status: "in-progress"}
			placeholderJSON, _ := json.Marshal(placeholder)

			p := &api.KVPair{Key: keyPath, Value: placeholderJSON, CreateIndex: 0}
			success, _, err := kv.CAS(p, nil) // Check and set with CreateIndex 0 is "put-if-absent"
			if err != nil {
				http.Error(w, "Failed to write to Consul", http.StatusInternalServerError)
				return
			}
			if !success {
				http.Error(w, "A concurrent request with the same idempotency key is in progress.", http.StatusConflict)
				return
			}

			// This function will be scheduled to be executed later in the function after the request is handled in the handler component,
			// 1. If any type of error happened in the handler, r will not be nil and the logic inside the if block will be called,
			// 2. If everything went well, r will be nil and the code inside the if block won't be called at all
			defer func() {
				if r := recover(); r != nil {
					kv.Delete(keyPath, nil)
					panic(r) // re-panic
				}
			}()

			// Process the actual request and capture the response.
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// Store the final response in Consul.
			// Only update if the request was successful (2xx status code).
			// If it was a client/server error, we delete the key to allow a retry.
			if rec.Code >= 200 && rec.Code < 300 {
				finalRecord := models.IdempotencyRecord{
					Status:     "completed",
					StatusCode: rec.Code,
					Body:       rec.Body.String(),
				}
				finalJSON, _ := json.Marshal(finalRecord)

				finalPair := &api.KVPair{Key: keyPath, Value: finalJSON}
				if _, err := kv.Put(finalPair, nil); err != nil {
					log.Printf("ERROR: Failed to save final response for key '%s': %v", idempotencyKey, err)
				} else {
					log.Printf("Saved final response for key '%s'", idempotencyKey)
				}
			} else {
				// The request failed, so we should allow it to be retried.
				kv.Delete(keyPath, nil)
			}

			// Write the captured response back to the original response writer.
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			io.Copy(w, rec.Body)
		})
	}
}
