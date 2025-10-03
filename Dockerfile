# Stage 1 (Build)
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
	-ldflags="-s -w" \
	-v \
	-trimpath \
	-o config-service \
	main.go

# Stage 2 (Final)
FROM gcr.io/distroless/static:latest

WORKDIR /app

COPY --from=builder /app/config-service /usr/bin/
COPY --from=builder /app/swagger.yaml ./swagger.yaml

ENTRYPOINT ["/usr/bin/config-service"]

EXPOSE 8000
