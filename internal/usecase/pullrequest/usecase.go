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
	Update(ctx context.Context, pr domain.PullRequest) error
}

type Usecase interface {
	Create(c context.Context, input CreateInput) (CreateOutput, error)

	// Пометить PR как MERGED (идемпотентная операция)
	Merge(c context.Context, input MergeInput) (MergeOutput, error)

	// Переназначить конкретного ревьювера на другого из его команды
	// Тут дохера ошибок проверить надо
	Reassign(c context.Context, input ReassignInput) (ReassignOutput, error)
}

type usecase struct {
	tx   txmanager.TxManager
	repo PullRequestRepository
	log  logger.Logger
}

func New(tx txmanager.TxManager, repo PullRequestRepository, log logger.Logger) *usecase {
	return &usecase{
		tx:   tx,
		repo: repo,
		log:  log,
	}
}

// func (u *usecase) Create(ctx context.Context, input CreateInput) (CreateOutput, error) {
// 	log := u.log.With("method", "Create", "author_id", input.AuthorID)

// 	author, err := u.repo.GetByID(ctx, input.AuthorID)
// 	if err != nil {
// 		log.Error("failed to get author", "err", err)
// 		return CreateOutput{}, err // тут лучше завернуть в доменную ошибку
// 	}

// 	// Кандидаты — активные из команды автора, без самого автора
// 	exclude := []string{author.ID}
// 	candidates, err := u.repo.ListActiveByTeamExcept(ctx, author.TeamID, exclude, 2)
// 	if err != nil {
// 		log.Error("failed to list candidates", "err", err)
// 		return CreateOutput{}, err
// 	}

// 	reviewerIDs := make([]string, 0, 2)
// 	for _, c := range candidates {
// 		reviewerIDs = append(reviewerIDs, c.ID)
// 	}

// 	pr := domain.PullRequest{
// 		// ID, скорее всего, генерит БД (serial/uuid), можешь оставить пустым
// 		Title:       input.Title,
// 		AuthorID:    author.ID,
// 		Status:      domain.PRStatusOpen,
// 		ReviewerIDs: reviewerIDs,
// 	}

// 	created, err := u.prRepo.Create(ctx, pr)
// 	if err != nil {
// 		log.Error("failed to create pr", "err", err)
// 		return CreateOutput{}, err
// 	}

// 	return CreateOutput{
// 		ID:          created.ID,
// 		Title:       created.Title,
// 		AuthorID:    created.AuthorID,
// 		Status:      string(created.Status),
// 		ReviewerIDs: created.ReviewerIDs,
// 	}, nil
// }
