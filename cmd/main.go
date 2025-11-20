package main

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

    _ = r.Run(":8080")
}

// MAIN TODO:
// Нужно сделать отдельный usecase для каждой доменной сущности
// Чтобы разграничить его интерфейс и не получить god service
