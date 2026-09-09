package core_transport_http_server

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

func NewConfig() (Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR is required")
	}

	shutdownTimeout, err := time.ParseDuration(os.Getenv("HTTP_SHUTDOWN_TIMEOUT"))
	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_SHUTDOWN_TIMEOUT: %w", err)
	}

	return Config{
		Addr:            addr,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}