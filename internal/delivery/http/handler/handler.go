package handler

import (
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/user"
)

type Handler struct {
	user user.Usecase
	team team.Usecase
	pr   pullrequest.Usecase
}

func NewHandler(userUC user.Usecase, teamUC team.Usecase, prUC pullrequest.Usecase) *Handler {
	return &Handler{
		user: userUC,
		team: teamUC,
		pr:   prUC,
	}
}

func newErrorResponse(code api.ErrorResponseErrorCode, msg string) api.ErrorResponse {
	var e api.ErrorResponse
	e.Error.Code = code
	e.Error.Message = msg
	return e
}

func mapErrorToErrorResponse(err error) api.ErrorResponse {
	// TODO
	return api.ErrorResponse{}
}

/*
func (h *Handler) PostTeamAdd(c *gin.Context) {
	// маппишь в domain-модель или сразу кидаешь в usecase
	team, err := h.teamUC.CreateOrUpdate(c, body)
	if err != nil {
		// маппишь доменную ошибку в ErrorResponse
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": team, // можно вернуть тот же Team или DTO
	})
}

func (h *Handler) GetTeamGet(c *gin.Context, params GetTeamGetParams) {
	team, err := h.teamUC.GetByName(c, params.TeamName)
	if err != nil {
		// если не найдено
		c.JSON(http.StatusNotFound, ErrorResponse{
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

	c.JSON(http.StatusOK, team)
}
*/
