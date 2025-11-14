package server

import (
	"context"
	"net/http"

	"github.com/pvj08/avito-autumn-2025/pkg/config"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type Server struct {
	logger logger.Logger
	config config.Server
	server *http.Server
}

func NewServer(logger logger.Logger, cfg config.Server, handler http.Handler) *Server {
	s := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		logger: logger,
		config: cfg,
		server: s,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	// стартуем в горутине
	go func() {
		s.logger.Info("starting HTTP server", "addr", s.config.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("http server error", "error", err)
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		// получен сигнал на завершение — делаем graceful shutdown
		s.logger.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()
		return s.server.Shutdown(shutdownCtx) // закрывает слушатель и корректно ждёт активные запросы
	case err := <-errCh:
		return err
	}
}
