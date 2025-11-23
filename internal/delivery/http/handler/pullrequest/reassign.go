package pullrequest

func (h *Handler) PostPullRequestReassign(c *gin.Context) {
	var req api.PostPullRequestReassignJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}
}
