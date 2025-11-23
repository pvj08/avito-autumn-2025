package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
)

func (h *Handler) PostTeamAdd(c *gin.Context) {
	var req api.PostTeamAddJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}

	out, err := h.team.Add(c.Request.Context(), team.AddInput{
		Team: team.Team{
			TeamName: req.TeamName,
			Members:  toDomainMembers(req.Members),
		},
	})

	if err != nil {
		errResp := mapErrorToErrorResponse(err)
		c.JSON(500, errResp) // TODO: proper status code
		return
	}

	resp := PostTeamAddResponse{
		Team: Team{
			TeamName: out.TeamName,
			Members:  toResponseMembers(out.Members),
		},
	}

	c.JSON(http.StatusOK, resp)
}
