package pullrequest

import (
	"math/rand"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

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
