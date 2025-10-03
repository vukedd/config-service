# Stage 1 (Build)
FROM golang:1.24.6-alpine AS builder

WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/

RUN CGO_ENABLED=0 go build \
	-ldflags="-s -w" \
	-v \
	-trimpath \
	-o config-service \
	main.go

RUN echo "ID=\"distroless\"" > /etc/os-release

# Stage 2 (Final)
FROM gcr.io/distroless/static:latest
COPY --from=builder /etc/os-release /etc/os-release

COPY --from=builder /app/config-service /usr/bin/

ENTRYPOINT ["/usr/bin/config-service"]

EXPOSE 8000
