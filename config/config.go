package config

import "os"

type Config struct {
	Address       string
	Port          string
	ConsulAddress string
	JaegerAddress string
}

func GetConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	consulAddress := os.Getenv("CONSUL_ADDRESS")
	if consulAddress == "" {
		consulAddress = "localhost:8500"
	}

	jaegerAddress := os.Getenv("JAEGER_ADDRESS")
	if jaegerAddress == "" {
		jaegerAddress = "http://localhost:14268/api/traces"
	}

	return Config{
		Address:       ":" + port,
		Port:          port,
		ConsulAddress: consulAddress,
		JaegerAddress: jaegerAddress,
	}
}