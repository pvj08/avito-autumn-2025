package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/handler"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/migrations"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/postgres"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/user"
	"github.com/pvj08/avito-autumn-2025/pkg/config"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
	httpserver "github.com/pvj08/avito-autumn-2025/pkg/server"
)

func main() {
	// Конфиг
	cfg := config.MustLoad()

	// Логгер
	log := logger.SetupLogger(cfg.Logger)

	// Контекст, завязанный на сигналы ОС
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	if err := AppRun(ctx, cfg, log); err != nil {
		log.Error("application run failed", "error", err)
		os.Exit(1)
	}
}

func AppRun(ctx context.Context, cfg config.Config, log logger.Logger) error {
	// Postgres
	pg, err := postgres.New(ctx, cfg.Postgres, log)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer func() {
		if err := pg.Close(); err != nil {
			log.Error("postgres close failed", "error", err)
		}
	}()

	if err := migrations.Up(pg.DB, cfg.Migrations); err != nil {
		return fmt.Errorf("migrations up: %w", err)
	}

	// TxManager
	txMgr := txmanager.NewSqlx(pg)

	// Repositories
	userRepo := postgres.NewUserRepo(pg)
	teamRepo := postgres.NewTeamRepo(pg)
	prRepo := postgres.NewPullRequestRepo(pg)

	// Usecases
	userUC := user.New(txMgr, userRepo, log)
	teamUC := team.New(txMgr, teamRepo, userRepo, log)
	prUC := pullrequest.New(txMgr, prRepo, teamRepo, log)

	// HTTP router
	r := gin.Default()

	// HTTP handlers
	h := handler.NewHandler(userUC, teamUC, prUC)

	// Регистрация хендлеров, сгенерированных из OpenAPI
	api.RegisterHandlers(r, h)

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// HTTP сервер-обёртка
	server := httpserver.NewServer(log, cfg.Server, r)

	// Горутинa, которая ждёт отмены контекста
	go func() {
		<-ctx.Done()

		log.Info("shutdown: context canceled")

		if err := server.Close(); err != nil {
			log.Error("server close failed", "err", err)
		} else {
			log.Info("server closed gracefully")
		}
	}()

	log.Info("starting http server", "addr", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))

	// Блокирующий запуск сервера
	if err := server.Run(); err != nil {
		// Ожидаемое завершение — http.ErrServerClosed
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server run: %w", err)
		}
	}

	log.Info("http server stopped")
	return nil
}
