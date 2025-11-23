package pullrequest

import (
	"context"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	GetByID(ctx context.Context, id string) (domain.PullRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]domain.PullRequest, error)
	// Update(ctx context.Context, pr domain.PullRequest) error
	Save(ctx context.Context, pr domain.PullRequest) error
}

type TeamReader interface {
	GetByUserID(ctx context.Context, userID string) (domain.Team, error)
}

type Usecase interface {
	Create(c context.Context, input CreateInput) (CreateOutput, error)
	// Пометить PR как MERGED (идемпотентная операция)
	Merge(c context.Context, input MergeInput) (MergeOutput, error)
	// Переназначить конкретного ревьювера на другого из его команды
	Reassign(c context.Context, input ReassignInput) (ReassignOutput, error)
	GetReview(c context.Context, input GetReviewInput) (GetReviewOutput, error)
}

type usecase struct {
	tx         txmanager.TxManager
	prRepo     PullRequestRepository
	teamReader TeamReader
	log        logger.Logger
}

func New(
	tx txmanager.TxManager,
	pr PullRequestRepository,
	team TeamReader,
	log logger.Logger,
) *usecase {
	return &usecase{
		tx:         tx,
		prRepo:     pr,
		teamReader: team,
		log:        log,
	}
}
