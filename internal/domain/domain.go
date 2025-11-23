package domain

import "time"

type PullRequestStatus string

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequest struct {
	// AssignedReviewers user_id назначенных ревьюверов (0..2)
	AssignedReviewers []string
	AuthorID          string
	PullRequestID     string
	PullRequestName   string
	Status            PullRequestStatus

	CreatedAt *time.Time
	MergedAt  *time.Time
}

type User struct {
	IsActive bool
	TeamName string
	UserID   string
	Username string
}

type Team struct {
	Members  []TeamMember
	TeamName string
}

type TeamMember struct {
	IsActive bool
	UserID   string
	Username string
}

func (pr *PullRequest) Merge() {
	pr.Status = PullRequestStatusMERGED
	now := time.Now()
	pr.MergedAt = &now
}
