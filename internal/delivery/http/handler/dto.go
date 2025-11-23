package handler

import (
	"time"

	"github.com/pvj08/avito-autumn-2025/internal/delivery/http/api"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/pullrequest"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
)

type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	MergedAt  *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	AuthorID        string `json:"author_id"`
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	Status          string `json:"status"`
}

type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type Team struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type PostPullRequestCreateResponse struct {
	PR PullRequest `json:"pr"`
}

type PostPullRequestMergeResponse struct {
	PR PullRequest `json:"pr"`
}

type PostPullRequestReassignResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

type PostTeamAddResponse struct {
	Team Team `json:"team"`
}

type GetTeamGetResponse struct {
	Team Team `json:"team"`
}

type GetUsersGetReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type PostUsersSetIsActiveResponse struct {
	U User `json:"user"`
}

func toDomainMembers(ms []api.TeamMember) []team.TeamMember {
	res := make([]team.TeamMember, len(ms))
	for i, m := range ms {
		res[i] = team.TeamMember{
			IsActive: m.IsActive,
			UserID:   m.UserId,
			Username: m.Username,
		}
	}
	return res
}

func toResponseMembers(members []team.TeamMember) []TeamMember {
	res := make([]TeamMember, len(members))
	for i, m := range members {
		res[i] = TeamMember{
			IsActive: m.IsActive,
			UserID:   m.UserID,
			Username: m.Username,
		}
	}
	return res
}

func toResponseReviews(prs []pullrequest.PullRequestShort) []PullRequestShort {
	res := make([]PullRequestShort, len(prs))
	for i, p := range prs {
		res[i] = PullRequestShort{
			AuthorID:        p.AuthorID,
			PullRequestID:   p.PullRequestID,
			PullRequestName: p.PullRequestName,
			Status:          string(p.Status),
		}
	}
	return res
}
