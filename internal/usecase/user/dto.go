package user

import "github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"

type User struct {
	IsActive bool
	TeamName string
	UserID   string
	Username string
}

type SetIsActiveInput struct {
	IsActive bool
	UserID   string
}

type SetIsActiveOutput struct {
	User
}

type GetReviewInput struct {
	UserID string
}

type GetReviewOutput struct {
	UserID       string
	PullRequests []pullrequest.PullRequestShort
}
