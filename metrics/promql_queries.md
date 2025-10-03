# Prometheus Queries za Configuration Service Metrike

## 1. Ukupan broj zahteva za prethodna 24 sata

```promql
# Ukupan broj zahteva u poslednja 24 sata
increase(config_service_http_requests_total[24h])

# Ukupan broj zahteva u poslednja 24 sata (grupisano)
sum(increase(config_service_http_requests_total[24h]))

# Ukupan broj zahteva po endpoint-u u poslednja 24 sata
sum by (endpoint) (increase(config_service_http_requests_total[24h]))

# Ukupan broj zahteva po metodi u poslednja 24 sata
sum by (method) (increase(config_service_http_requests_total[24h]))
```

## 2. Broj uspešnih zahteva (status kodovi 2xx, 3xx) za prethodna 24 sata

```promql
# Uspešni zahtevi u poslednja 24 sata
sum(increase(config_service_http_requests_successful_total[24h]))

# Uspešni zahtevi po endpoint-u u poslednja 24 sata
sum by (endpoint) (increase(config_service_http_requests_successful_total[24h]))

# Uspešni zahtevi po metodi u poslednja 24 sata
sum by (method) (increase(config_service_http_requests_successful_total[24h]))

# Procenat uspešnih zahteva u poslednja 24 sata
(sum(increase(config_service_http_requests_successful_total[24h])) / sum(increase(config_service_http_requests_total[24h]))) * 100
```

## 3. Broj neuspešnih zahteva (status kodovi 4xx, 5xx) za prethodna 24 sata

```promql
# Neuspešni zahtevi u poslednja 24 sata
sum(increase(config_service_http_requests_failed_total[24h]))

# Neuspešni zahtevi po endpoint-u u poslednja 24 sata
sum by (endpoint) (increase(config_service_http_requests_failed_total[24h]))

# Neuspešni zahtevi po status klasi u poslednja 24 sata
sum by (status_class) (increase(config_service_http_requests_failed_total[24h]))

# 4xx errori u poslednja 24 sata
sum(increase(config_service_http_requests_failed_total{status_class="4xx"}[24h]))

# 5xx errori u poslednja 24 sata
sum(increase(config_service_http_requests_failed_total{status_class="5xx"}[24h]))

# Procenat neuspešnih zahteva u poslednja 24 sata
(sum(increase(config_service_http_requests_failed_total[24h])) / sum(increase(config_service_http_requests_total[24h]))) * 100
```

## 4. Prosečno vreme izvršavanja zahteva za svaki endpoint

```promql
# Prosečno vreme odgovora za sve endpoint-e u poslednja 24 sata
rate(config_service_http_response_duration_seconds_sum[24h]) / rate(config_service_http_response_duration_seconds_count[24h])

# Prosečno vreme odgovora po endpoint-u u poslednja 24 sata
rate(config_service_http_response_duration_seconds_sum[24h]) / rate(config_service_http_response_duration_seconds_count[24h]) by (endpoint)

# Prosečno vreme odgovora po endpoint-u i metodi u poslednja 24 sata
rate(config_service_http_response_duration_seconds_sum[24h]) / rate(config_service_http_response_duration_seconds_count[24h]) by (endpoint, method)

# 95. percentil vremena odgovora po endpoint-u
histogram_quantile(0.95, rate(config_service_http_response_duration_seconds_bucket[5m])) by (endpoint)

# 99. percentil vremena odgovora po endpoint-u
histogram_quantile(0.99, rate(config_service_http_response_duration_seconds_bucket[5m])) by (endpoint)

# Medijan vremena odgovora po endpoint-u
histogram_quantile(0.5, rate(config_service_http_response_duration_seconds_bucket[5m])) by (endpoint)
```

## 5. Broj zahteva u jedinici vremena za svaki endpoint za prethodna 24 sata

```promql
# Broj zahteva po minuti za svaki endpoint (poslednja 24 sata)
rate(config_service_http_request_rate_per_minute[24h]) * 60

# Broj zahteva po sekundi za svaki endpoint (poslednja 24 sata)
rate(config_service_http_requests_total[24h])

# Broj zahteva po minuti po endpoint-u (poslednja 24 sata)
sum by (endpoint) (rate(config_service_http_request_rate_per_minute[24h])) * 60

# Broj zahteva po sekundi po endpoint-u (poslednja 24 sata)
sum by (endpoint) (rate(config_service_http_requests_total[24h]))

# Broj zahteva po minuti po endpoint-u i metodi (poslednja 24 sata)
sum by (endpoint, method) (rate(config_service_http_request_rate_per_minute[24h])) * 60

# Broj zahteva po sekundi po endpoint-u i metodi (poslednja 24 sata)
sum by (endpoint, method) (rate(config_service_http_requests_total[24h]))
```

## 6. Dodatne korisne metrike

```promql
# Trenutno aktivni zahtevi
config_service_http_requests_in_flight

# Najviši broj aktivnih zahteva u poslednja 24 sata
max_over_time(config_service_http_requests_in_flight[24h])

# Prosečan broj aktivnih zahteva u poslednja 24 sata
avg_over_time(config_service_http_requests_in_flight[24h])

# Top 5 najsporijih endpoint-a po prosečnom vremenu odgovora
topk(5, rate(config_service_http_response_duration_seconds_sum[24h]) / rate(config_service_http_response_duration_seconds_count[24h]) by (endpoint))

# Top 5 endpoint-a sa najviše zahteva
topk(5, sum by (endpoint) (increase(config_service_http_requests_total[24h])))

# Top 5 endpoint-a sa najviše grešaka
topk(5, sum by (endpoint) (increase(config_service_http_requests_failed_total[24h])))
```

## 7. Alerting pravila (za Prometheus Alert Manager)

```promql
# High error rate (više od 5% grešaka u poslednja 5 minuta)
(sum(rate(config_service_http_requests_failed_total[5m])) / sum(rate(config_service_http_requests_total[5m]))) * 100 > 5

# High response time (95. percentil veći od 1 sekunde)
histogram_quantile(0.95, rate(config_service_http_response_duration_seconds_bucket[5m])) > 1

# High request rate (više od 100 zahteva po sekundi)
sum(rate(config_service_http_requests_total[1m])) > 100

# Service down (nema zahteva u poslednja 2 minuta)
sum(rate(config_service_http_requests_total[2m])) == 0
```

## 9. Korisni upiti za monitoring

```promql
# Trenutni QPS (queries per second)
sum(rate(config_service_http_requests_total[1m]))

# Trenutni error rate (procenat)
sum(rate(config_service_http_requests_failed_total[1m])) / sum(rate(config_service_http_requests_total[1m])) * 100

# Trenutno vreme odgovora (95. percentil)
histogram_quantile(0.95, rate(config_service_http_response_duration_seconds_bucket[5m]))

# Najaktivniji endpoint-i
topk(10, sum by (endpoint) (rate(config_service_http_requests_total[1h])))

# Endpoint-i sa najvećim brojem grešaka
topk(10, sum by (endpoint) (rate(config_service_http_requests_failed_total[1h])))
```
