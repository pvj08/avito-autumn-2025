package handler

import "github.com/pvj08/avito-autumn-2025/internal/delivery/http/handler"

func (h *handler.Handler) PostPullRequestReassign(c *gin.Context) {
	var req api.PostPullRequestReassignJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}
}
