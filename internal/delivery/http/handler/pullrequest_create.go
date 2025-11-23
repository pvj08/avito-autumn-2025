package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
)

func (h *Handler) PostPullRequestCreate(c *gin.Context) {
	var req api.PostPullRequestCreateJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}

	out, err := h.pr.Create(c.Request.Context(), pullrequest.CreateInput{
		AuthorID:        req.AuthorId,
		PullRequestID:   req.PullRequestId,
		PullRequestName: req.PullRequestName,
	})

	if err != nil {
		errResp := mapErrorToErrorResponse(err)
		c.JSON(500, errResp) // TODO: proper status code
		return
	}

	resp := PostPullRequestCreateResponse{
		PR: PullRequest{
			PullRequestID:     out.PullRequestID,
			PullRequestName:   out.PullRequestName,
			AuthorID:          out.AuthorID,
			Status:            string(out.Status),
			AssignedReviewers: out.AssignedReviewers,
		},
	}

	c.JSON(http.StatusOK, resp)
}
