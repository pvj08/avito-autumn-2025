package pullrequest

import (
	"time"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

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

type GetReviewInput struct {
	UserID string
}

type GetReviewOutput struct {
	UserID       string
	PullRequests []PullRequestShort
}

// toDomainPullRequest маппит входной CreateInput + выбранных ревьюверов в domain.PullRequest.
func fromAddInputToDomainPR(input CreateInput, reviewers []string) domain.PullRequest {
	return domain.PullRequest{
		PullRequestID:     input.PullRequestID,
		PullRequestName:   input.PullRequestName,
		AuthorID:          input.AuthorID,
		AssignedReviewers: reviewers,
		Status:            domain.PullRequestStatusOPEN,
	}
}

// fromDomainPullRequest маппит domain.PullRequest в DTO PullRequest (ответ API).
func fromDomainPullRequest(pr domain.PullRequest) PullRequest {
	return PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		AssignedReviewers: pr.AssignedReviewers,
		Status:            PullRequestStatus(pr.Status),
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func fromDomainPRtoShortPR(pr domain.PullRequest) PullRequestShort {
	return PullRequestShort{
		AuthorID:        pr.AuthorID,
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		Status:          PullRequestShortStatus(pr.Status),
	}
}

func fromDomainPRtoShortPRSlice(prs []domain.PullRequest) []PullRequestShort {
	if len(prs) == 0 {
		return []PullRequestShort{}
	}

	out := make([]PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		out = append(out, fromDomainPRtoShortPR(pr))
	}

	return out
}
