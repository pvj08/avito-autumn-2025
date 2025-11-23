package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/handler"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/postgres"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/user"
	"github.com/pvj08/avito-autumn-2025/pkg/config"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
	httpserver "github.com/pvj08/avito-autumn-2025/pkg/server"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Logger)

	err := AppRun(context.Background(), cfg, log)
	if err != nil {
		log.Error("application run failed", "error", err)
		panic(err)
	}
}

func AppRun(ctx context.Context, cfg config.Config, log logger.Logger) error {
	// Postgres
	pg, err := postgres.New(ctx, cfg.Postgres, log)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	// Repos
	userRepo := postgres.NewUserRepo(pg)
	teamRepo := postgres.NewTeamRepo(pg)
	prRepo := postgres.NewPullRequestRepo(pg)

	// Usecase (Service)
	userUC := user.New(log, userRepo)
	teamUC := team.New(log, teamRepo)
	prUC := pullrequest.New(log, prRepo)

	r := gin.Default()
	h := handler.NewHandler(userUC, teamUC, prUC)
	server := httpserver.NewServer(log, cfg.Server, r)

	// функция из openapi_server.gen.go
	api.RegisterHandlers(r, h)

	// опционально: healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Приложение запущено и готово к работе
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig // ждём здесь сигнала (Ctrl+C или SIGTERM)
	log.Info("shutdown signal received")

	// Закрываем ресурсы
	httpServer.Close()
	pgPool.Close()

	return nil
}

// MAIN TODO:
// Нужно сделать отдельный usecase для каждой доменной сущности
// Чтобы разграничить его интерфейс и не получить god service

/*

Небольшая рефлексия.
Кодогенерацией сделан рутинг и дто модели.
Что мне нужно сделать? Очевидно, хендлеры и бизнес логику.

в models
*/
