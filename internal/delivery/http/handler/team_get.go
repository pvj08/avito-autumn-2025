package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
)

func (h *Handler) GetTeamGet(c *gin.Context, params api.GetTeamGetParams) {
	out, err := h.team.Get(c.Request.Context(), team.GetInput{
		TeamName: params.TeamName,
	})

	if err != nil {
		status, resp := mapDomainErrorToErrorResponse(err)
		if status == http.StatusInternalServerError {
			h.log.Error("internal error", "error", err)
		}
		c.JSON(status, resp)
		return
	}

	resp := GetTeamGetResponse{
		Team: Team{
			TeamName: out.TeamName,
			Members:  toResponseMembers(out.Members),
		},
	}

	c.JSON(http.StatusOK, resp)
}
