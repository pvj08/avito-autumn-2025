package handler

func (h *Handler) PostPullRequestCreate(c *gin.Context) {
	var req api.PostPullRequestCreateJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}

	output, err := h.profileService.CreateProfile(r.Context(), PullRequest.CreateInput{ // TODO: DTO struct
		req.AuthorId,
		req.PullRequestId,
		req.PullRequestName
	})

	if err != nil {
		errResp := mapErrorToErrorResponse(err) // TODO map func
		c.JSON(errResp.HTTPStatus, errResp)
		return
	}

	c.JSON(http.StatusOK, output) // TODO: map output to delivery DTO
	/*
	или мне тут похуй и я могу не мапить, а просто возвращать доменный объект?
	*/
}