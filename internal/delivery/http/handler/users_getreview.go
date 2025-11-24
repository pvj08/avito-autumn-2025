package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
)

func (h *Handler) GetUsersGetReview(c *gin.Context, params api.GetUsersGetReviewParams) {
	out, err := h.pr.GetReview(c.Request.Context(), pullrequest.GetReviewInput{
		UserID: params.UserId,
	})

	if err != nil {
		status, resp := mapDomainErrorToErrorResponse(err)
		if status == http.StatusInternalServerError {
			h.log.Error("internal error", "error", err)
		}
		c.JSON(status, resp)
		return
	}

	resp := GetUsersGetReviewResponse{
		UserID:       out.UserID,
		PullRequests: toResponseReviews(out.PullRequests),
	}

	c.JSON(http.StatusOK, resp)
}
