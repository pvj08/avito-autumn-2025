package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
)

// Пометить PR как MERGED (идемпотентная операция)
func (h *Handler) PostPullRequestMerge(c *gin.Context) {
	var req api.PostPullRequestMergeJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("BAD_REQUEST", "invalid JSON"))
		return
	}

	out, err := h.pr.Merge(c.Request.Context(), pullrequest.MergeInput{
		PullRequestID: req.PullRequestId,
	})

	if err != nil {
		status, resp := mapDomainErrorToErrorResponse(err)
		if status == http.StatusInternalServerError {
			h.log.Error("internal error", "error", err)
		}
		c.JSON(status, resp)
		return
	}

	resp := PostPullRequestMergeResponse{
		PR: PullRequest{
			PullRequestID:     out.PullRequestID,
			PullRequestName:   out.PullRequestName,
			AuthorID:          out.AuthorID,
			Status:            string(out.Status),
			AssignedReviewers: out.AssignedReviewers,
			MergedAt:          out.MergedAt,
		},
	}

	c.JSON(http.StatusOK, resp)
}
