package handler

func (h *Handler) PostPullRequestMerge(c *gin.Context) {
	var req api.PostPullRequestMergeJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}
}
