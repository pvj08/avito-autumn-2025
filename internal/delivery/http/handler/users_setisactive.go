package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/user"
)

func (h *Handler) PostUsersSetIsActive(c *gin.Context) {
	var req api.PostUsersSetIsActiveJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}

	out, err := h.user.SetIsActive(c.Request.Context(), user.SetIsActiveInput{
		UserID:   req.UserId,
		IsActive: req.IsActive,
	})

	if err != nil {
		status, resp := mapDomainErrorToErrorResponse(err)
		if status == http.StatusInternalServerError {
			h.log.Error("internal error", "error", err)
		}
		c.JSON(status, resp)
		return
	}

	resp := PostUsersSetIsActiveResponse{
		U: User{
			UserID:   out.UserID,
			Username: out.Username,
			IsActive: out.IsActive,
			TeamName: out.TeamName,
		},
	}

	c.JSON(http.StatusOK, resp)
}
