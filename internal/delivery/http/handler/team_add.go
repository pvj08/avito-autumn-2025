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
		status, resp := mapDomainErrorToErrorResponse(err)
		if status == http.StatusInternalServerError {
			h.log.Error("internal error", "error", err)
		}
		c.JSON(status, resp)
		return
	}

	resp := PostTeamAddResponse{
		Team: Team{
			TeamName: out.TeamName,
			Members:  toResponseMembers(out.Members),
		},
	}

	c.JSON(http.StatusCreated, resp)
}
