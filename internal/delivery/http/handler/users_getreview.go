package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/user"
)

func (h *Handler) GetUsersGetReview(c *gin.Context, params api.GetUsersGetReviewParams) {
	out, err := h.user.GetReview(c.Request.Context(), user.GetReviewInput{
		UserID: params.UserId,
	})

	if err != nil {
		errResp := mapErrorToErrorResponse(err)
		// TODO: различать 404 / 500 и т.п.
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	resp := GetUsersGetReviewResponse{
		UserID:       out.UserID,
		PullRequests: toResponseReviews(out.PullRequests),
	}

	c.JSON(http.StatusOK, resp)
}
