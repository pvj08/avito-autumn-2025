package user

import "github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"

type User struct {
	IsActive bool
	TeamName string
	UserId   string
	Username string
}

type SetIsActiveInput struct {
	IsActive bool
	UserId   string
}

type SetIsActiveOutput struct {
	User
}

type GetReviewInput struct {
	UserId string
}

type GetReviewOutput struct {
	UserId       string
	PullRequests pullrequest.PullRequestShort
}
