package core_transport_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(config Config, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    config.Addr,
			Handler: handler,
		},
		shutdownTimeout: config.ShutdownTimeout,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		err := s.server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listen and serve HTTP: %w", err)

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.shutdownTimeout,
		)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			_ = s.server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		return nil
	}
}
