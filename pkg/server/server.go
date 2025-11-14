package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	delivery "github.com/pvj08/effective-mobile/internal/delivery/http"
	"github.com/pvj08/effective-mobile/pkg/config"
	"github.com/pvj08/effective-mobile/pkg/logger"
)

type Server struct {
	logger logger.Logger
	config config.Server
	engine *gin.Engine
	server *http.Server
}

func NewServer(logger logger.Logger, cfg config.Server, deps delivery.Deps) *Server {
	eng := delivery.NewRouter(deps)

	s := &http.Server{
		Addr:         cfg.Addr,
		Handler:      eng,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		logger: logger,
		config: cfg,
		engine: eng,
		server: s,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	// стартуем в горутине
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		// получен сигнал на завершение — делаем graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()
		return s.server.Shutdown(shutdownCtx) // закрывает слушатель и корректно ждёт активные запросы
	case err := <-errCh:
		return err
	}
}
