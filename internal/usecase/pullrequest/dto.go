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
	AuthorId          string
	CreatedAt         *time.Time
	MergedAt          *time.Time
	PullRequestId     string
	PullRequestName   string
	Status            PullRequestStatus
}

// PullRequestShort defines model for PullRequestShort.
type PullRequestShort struct {
	AuthorId        string
	PullRequestId   string
	PullRequestName string
	Status          PullRequestShortStatus
}

type CreateInput struct {
	AuthorId        string
	PullRequestId   string
	PullRequestName string
}

type CreateOutput struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            PullRequestStatus
	AssignedReviewers []string
}

type MergeInput struct {
	PullRequestId string
}

type MergeOutput struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            PullRequestStatus
	AssignedReviewers []string
	MergedAt          *time.Time
}

type ReassignInput struct {
	PullRequestId string
	UserID        string
}

type ReassignOutput struct {
	PR         PullRequest
	ReplacedBy string
}
