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
