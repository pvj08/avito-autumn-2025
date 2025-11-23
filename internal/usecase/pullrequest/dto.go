package pullrequest

import "time"

type PullRequestStatus string

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequestShortStatus string

const (
	PullRequestShortStatusMERGED PullRequestShortStatus = "MERGED"
	PullRequestShortStatusOPEN   PullRequestShortStatus = "OPEN"
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

// PullRequestShort defines model for PullRequestShort.
type PullRequestShort struct {
	AuthorID        string
	PullRequestID   string
	PullRequestName string
	Status          PullRequestShortStatus
}

type CreateInput struct {
	AuthorID        string
	PullRequestID   string
	PullRequestName string
}

type CreateOutput struct {
	PullRequest
}

type MergeInput struct {
	PullRequestID string
}

type MergeOutput struct {
	PullRequest
}

type ReassignInput struct {
	PullRequestID string
	UserID        string
}

type ReassignOutput struct {
	PullRequest
	ReplacedBy string
}
