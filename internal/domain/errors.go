package domain

import "errors"

var (
	// Эти ошибки надо было в репозиторий вынести
	ErrAlreadyExists = errors.New("entity already exists")
	ErrNotFound      = errors.New("entity not found")

	// Бизнес ошибки
	ErrTeamExists = errors.New("team_name already exists")
	ErrUserExists = errors.New("user_id already exists") // 409 conflict
	ErrPrExists   = errors.New("PR id already exists")

	ErrMerged      = errors.New("cannot reassign on merged PR")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
)
