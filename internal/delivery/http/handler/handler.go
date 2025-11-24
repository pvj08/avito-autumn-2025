package handler

import (
	"errors"
	"net/http"

	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/user"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type Handler struct {
	user user.Usecase
	team team.Usecase
	pr   pullrequest.Usecase
	log  logger.Logger
}

func NewHandler(userUC user.Usecase, teamUC team.Usecase, prUC pullrequest.Usecase, l logger.Logger) *Handler {
	return &Handler{
		user: userUC,
		team: teamUC,
		pr:   prUC,
		log:  l,
	}
}

func newErrorResponse(code api.ErrorResponseErrorCode, msg string) api.ErrorResponse {
	var e api.ErrorResponse
	e.Error.Code = code
	e.Error.Message = msg
	return e
}

func mapDomainErrorToErrorResponse(err error) (int, api.ErrorResponse) {
	var out api.ErrorResponse

	var status int
	var code api.ErrorResponseErrorCode
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		code = api.NOTFOUND

	case errors.Is(err, domain.ErrTeamExists):
		status = http.StatusBadRequest
		code = api.TEAMEXISTS

	case errors.Is(err, domain.ErrUserExists):
		status = http.StatusConflict
		code = api.ErrorResponseErrorCode("USER EXISTS")

	case errors.Is(err, domain.ErrPrExists):
		status = http.StatusConflict
		code = api.PREXISTS

	case errors.Is(err, domain.ErrMerged):
		status = http.StatusConflict
		code = api.PRMERGED

	case errors.Is(err, domain.ErrNotAssigned):
		status = http.StatusConflict
		code = api.NOTASSIGNED

	case errors.Is(err, domain.ErrNoCandidate):
		status = http.StatusConflict
		code = api.NOCANDIDATE

	default:
		// fallback: внутренняя ошибка
		status = http.StatusInternalServerError
		code = api.ErrorResponseErrorCode("INTERNAL_ERROR")
		err = errors.New("something went wrong")
	}

	out.Error.Message = err.Error()
	out.Error.Code = code

	return status, out
}
