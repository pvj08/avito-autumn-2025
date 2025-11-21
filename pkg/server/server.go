package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type Config struct {
	Host            string        `env:"HTTP_HOST"             env-default:"0.0.0.0"`
	Port            int           `env:"HTTP_PORT"             env-default:"8080"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT"     env-default:"5s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT"    env-default:"10s"`
	IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT"     env-default:"60s"`
	MaxHeaderBytes  int           `env:"HTTP_MAX_HEADER_BYTES" env-default:"1048576"` // 1MB
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(logger logger.Logger, c Config, handler http.Handler) *Server {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	s := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		IdleTimeout:    c.IdleTimeout,
		MaxHeaderBytes: c.MaxHeaderBytes,
	}

	return &Server{
		server:          s,
		shutdownTimeout: c.ShutdownTimeout,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
