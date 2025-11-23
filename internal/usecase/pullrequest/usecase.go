package pullrequest

import (
	"context"
	"math/rand"
	"time"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/internal/usecase/team"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Usecase interface {
	Create(c context.Context, input CreateInput) (CreateOutput, error)

	// Пометить PR как MERGED (идемпотентная операция)
	Merge(c context.Context, input MergeInput) (MergeOutput, error)

	// Переназначить конкретного ревьювера на другого из его команды
	// Тут дохера ошибок проверить надо
	Reassign(c context.Context, input ReassignInput) (ReassignOutput, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	GetByID(ctx context.Context, id string) (domain.PullRequest, error)
	Update(ctx context.Context, pr domain.PullRequest) error
	Save(ctx context.Context, pr domain.PullRequest) error
}

type usecase struct {
	tx       txmanager.TxManager
	prRepo   PullRequestRepository
	teamRepo team.TeamRepository
	log      logger.Logger
}

func New(
	tx txmanager.TxManager,
	prRepo PullRequestRepository,
	teamRepo team.TeamRepository,
	log logger.Logger,
) *usecase {
	return &usecase{
		tx:       tx,
		prRepo:   prRepo,
		teamRepo: teamRepo,
		log:      log,
	}
}

// выбирает до 2 активных ревьюверов из команды, исключая ID автора.
func chooseReviewers(team domain.Team, authorID string) ([]string, error) {
	candidates := make([]string, 0, len(team.Members))

	// 1. фильтруем подходящих
	for _, m := range team.Members {
		if !m.IsActive {
			continue
		}
		if m.UserID == authorID {
			continue
		}
		candidates = append(candidates, m.UserID)
	}

	// 2. нет кандидатов
	if len(candidates) == 0 {
		return []string{}, domain.ErrNoCandidate
	}

	// 3. один кандидат
	if len(candidates) == 1 {
		return []string{candidates[0]}, nil
	}

	// 4. два или больше — выбираем случайные два
	// перемешаем кандидатов
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	// берём первые два
	return candidates[:2], nil
}

// выбирает одного ревьювера из команды, исключая ID автора и exceptID.
func reassignChooseReviewer(team domain.Team, authorID, exceptID string) (string, error) {
	candidates := make([]string, 0, len(team.Members))

	// 1. фильтруем активных участников команды, исключая автора и exceptID
	for _, m := range team.Members {
		if !m.IsActive {
			continue
		}
		if m.UserID == authorID {
			continue
		}
		if m.UserID == exceptID {
			continue
		}
		candidates = append(candidates, m.UserID)
	}

	// нет доступных ревьюверов
	if len(candidates) == 0 {
		return "", domain.ErrNoCandidate
	}

	// выбираем случайно одного
	return candidates[rand.Intn(len(candidates))], nil
}
