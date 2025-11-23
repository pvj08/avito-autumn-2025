package domain

import "time"

type UserID string
type TeamID string
type PullRequestID string

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type User struct {
	ID       UserID
	Name     string
	IsActive bool

	// один пользователь — одна команда
	TeamID TeamID
}

type Team struct {
	ID   TeamID
	Name string

	MemberIDs []UserID
}

type PullRequest struct {
	ID    PullRequestID
	Title string

	AuthorID UserID
	Status   PRStatus

	ReviewerIDs []UserID

	MergedAt  *time.Time
	CreatedAt *time.Time
}
