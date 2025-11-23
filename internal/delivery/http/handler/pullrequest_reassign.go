package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
)

func (h *Handler) PostPullRequestReassign(c *gin.Context) {
	var req api.PostPullRequestReassignJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}

	out, err := h.pr.Reassign(c.Request.Context(), pullrequest.ReassignInput{
		PullRequestID: req.PullRequestId,
		UserID:        req.OldUserId,
	})

	if err != nil {
		errResp := mapErrorToErrorResponse(err)
		c.JSON(500, errResp) // TODO: proper status code
		return
	}

	resp := PostPullRequestReassignResponse{
		PR: PullRequest{
			PullRequestID:     out.PullRequestID,
			PullRequestName:   out.PullRequestName,
			AuthorID:          out.AuthorID,
			Status:            string(out.Status),
			AssignedReviewers: out.AssignedReviewers,
			MergedAt:          out.MergedAt,
		},
		ReplacedBy: out.ReplacedBy,
	}

	c.JSON(http.StatusOK, resp)
}
