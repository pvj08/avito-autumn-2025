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
		errResp := mapErrorToErrorResponse(err)
		// TODO: различать 404 / 500 и т.п.
		c.JSON(http.StatusInternalServerError, errResp)
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
