package handler

import "github.com/pvj08/avito-autumn-2025/pkg/logger"

type Handler struct {
	TeamUC TeamUsecase
	UserUC UserUsecase
	PrUC   PullRequestUsecase
}

func NewHandler(userUC user.Usecase, teamUC team.Usecase, prUC pullrequest.Usecase) *Handler {
    return &Handler{
        userUC: userUC,
        teamUC: teamUC,
        prUC:   prUC,
    }
}

/*
// ServerInterface represents all server handlers.
type ServerInterface interface {
	PostPullRequestCreate(c *gin.Context)
	PostPullRequestMerge(c *gin.Context)
	PostPullRequestReassign(c *gin.Context)
	PostTeamAdd(c *gin.Context)
	GetTeamGet(c *gin.Context, params GetTeamGetParams)
	GetUsersGetReview(c *gin.Context, params GetUsersGetReviewParams)
	PostUsersSetIsActive(c *gin.Context)
}
*/

func (h *Handler) PostTeamAdd(ctx *gin.Context) {
	var body Team // Team из openapi_models.gen.go // TODO: ⚠️❌⚠️❌⚠️PostTeamAddJSONRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "invalid JSON",
			},
		})
		return
	}

	// маппишь в domain-модель или сразу кидаешь в usecase
	team, err := h.teamUC.CreateOrUpdate(ctx, body)
	if err != nil {
		// маппишь доменную ошибку в ErrorResponse
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"team": team, // можно вернуть тот же Team или DTO
	})
}

func (h *Handler) GetTeamGet(ctx *gin.Context, params GetTeamGetParams) {
	team, err := h.teamUC.GetByName(ctx, params.TeamName)
	if err != nil {
		// если не найдено
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "NOT_FOUND",
				Message: "team not found",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, team)
}

func teamDTOToDomain(t Team) domain.Team { ... }
func prDTOToDomain(p PullRequest) domain.PullRequest { ... }