package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

    // собираешь зависимости (юзкейсы, репозитории и т.п.)
    teamUC := usecase.NewTeamUsecase(...)
    userUC := usecase.NewUserUsecase(...)
    prUC   := usecase.NewPullRequestUsecase(...)

    h := &http.Handler{
        teamUC: teamUC,
        userUC: userUC,
        prUC:   prUC,
    }

    // функция из openapi_server.gen.go
    RegisterHandlers(r, h)

    // опционально: healthcheck
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

	srv := server.NewServer(log, cfg.Server, router)

	// старт сервера в отдельной горутине
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", err)
		}
	}()

	// слушаем системные сигналы OS (Ctrl+C, docker stop, k8s SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // ждём первого сигнала

	log.Info("shutdown signal received")

	// здесь используется ShutdownTimeout
	if err := srv.Close(); err != nil {
		log.Error("graceful shutdown failed", err)
	}
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

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.LogLevel)


}