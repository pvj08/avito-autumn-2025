package handler

func (h *Handler) PostUsersSetIsActive(c *gin.Context) {
	var req api.PostUsersSetIsActiveJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}
}
