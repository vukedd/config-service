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
)

// IdempotencyMiddleware obezbeđuje idempotent operacije koristeći Consul kao storage
func IdempotencyMiddleware(consulClient *api.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idempotencyKey := r.Header.Get("Idempotency-Key")

			// Ako nema ključa, samo nastavi sa sledećim handler-om
			if idempotencyKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			kv := consulClient.KV()
			keyPath := fmt.Sprintf("idempotency/%s", idempotencyKey)

			// Proveri da li ključ već postoji
			pair, _, err := kv.Get(keyPath, nil)
			if err != nil {
				http.Error(w, "Failed to connect to Consul", http.StatusInternalServerError)
				return
			}

			if pair != nil {
				// Ključ postoji, proveri stanje
				var record models.IdempotencyRecord
				if err := json.Unmarshal(pair.Value, &record); err != nil {
					http.Error(w, "Failed to parse stored record", http.StatusInternalServerError)
					return
				}

				// Record je označen kao završen, vrati keširani odgovor (status i body)
				if record.Status == models.StatusCompleted {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(record.StatusCode)
					w.Write([]byte(record.Body))
					return
				}

				// Record je označen kao u toku, odbaci zahtev da se izbegne race condition
				if record.Status == models.StatusInProgress {
					http.Error(w, "Request with this idempotency key is already in progress.", http.StatusConflict)
					return
				}
			}

			// Ključ ne postoji. Kreiraj placeholder (lock)
			placeholder := &models.IdempotencyRecord{Status: models.StatusInProgress}
			placeholderJSON, _ := json.Marshal(placeholder)

			p := &api.KVPair{Key: keyPath, Value: placeholderJSON, CreateIndex: 0}
			success, _, err := kv.CAS(p, nil) // Check and set sa CreateIndex 0 je "put-if-absent"
			if err != nil {
				http.Error(w, "Failed to write to Consul", http.StatusInternalServerError)
				return
			}
			if !success {
				http.Error(w, "A concurrent request with the same idempotency key is in progress.", http.StatusConflict)
				return
			}

			// Ova funkcija će biti zakazana da se izvršava kasnije nakon što se zahtev obradi u handler komponenti
			// 1. Ako se dogodi bilo kakva greška u handler-u, r neće biti nil i logika unutar if bloka će biti pozvana
			// 2. Ako sve prođe dobro, r će biti nil i kod unutar if bloka se neće izvršiti
			defer func() {
				if r := recover(); r != nil {
					kv.Delete(keyPath, nil)
					panic(r) // re-panic
				}
			}()

			// Obradi stvarni zahtev i uhvati odgovor
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// Sačuvaj finalni odgovor u Consul
			// Ažuriraj samo ako je zahtev bio uspešan (2xx status kod)
			// Ako je bila client/server greška, obriši ključ da dozvoliš retry
			if rec.Code >= 200 && rec.Code < 300 {
				finalRecord := models.IdempotencyRecord{
					Status:     models.StatusCompleted,
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
				// Zahtev nije uspešan, obriši ključ da dozvoliš retry
				kv.Delete(keyPath, nil)
			}

			// Upiši uhvaćeni odgovor nazad u originalni response writer
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			io.Copy(w, rec.Body)
		})
	}
}