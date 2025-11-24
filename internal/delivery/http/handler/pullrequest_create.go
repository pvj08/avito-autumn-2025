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
		status, resp := mapDomainErrorToErrorResponse(err)
		if status == http.StatusInternalServerError {
			h.log.Error("internal error", "error", err)
		}
		c.JSON(status, resp)
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

	c.JSON(http.StatusCreated, resp)
}
