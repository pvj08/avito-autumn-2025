package api

import (
	"github.com/gin-gonic/gin" // или chi/echo — что любишь

	"github.com/pvj08/avito-autumn-2025/internal/usecase"
)

func NewRouter(prUC usecase.PullRequestUseCase, teamUC usecase.TeamUseCase, statsUC usecase.StatsUseCase) *gin.Engine {
	r := gin.Default()

	prHandler := NewPRHandler(prUC)
	teamHandler := NewTeamHandler(teamUC)
	statsHandler := NewStatsHandler(statsUC)

	api := r.Group("/api")

	// PR endpoints
	{
		api.POST("/prs", prHandler.Create)
		api.POST("/prs/:id/merge", prHandler.Merge)
		api.POST("/prs/:id/reassign", prHandler.ReassignReviewer)
		api.GET("/prs", prHandler.ListByReviewer) // ?reviewer_id=
	}

	// Team / Users
	{
		api.POST("/teams", teamHandler.CreateTeam)
		api.POST("/users", teamHandler.CreateUser)
		api.PATCH("/users/:id/active", teamHandler.SetUserActive)
		api.POST("/teams/:id/deactivate", teamHandler.BulkDeactivate) // доп-задание
	}

	// Stats
	{
		api.GET("/stats/assignments", statsHandler.ListAssignmentsByUser)
	}

	return r
}
