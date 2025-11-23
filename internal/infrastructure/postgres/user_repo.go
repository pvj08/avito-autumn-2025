package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
)

// type UserRepository interface {
// 	Create(ctx context.Context, user domain.User) error
// 	UpdateSetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
// }

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) exec(ctx context.Context) dbExecutor {
	if tx := txmanager.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}

// Create — создать пользователя.
// Если user_id уже существует — вернёт domain.ErrAlreadyExists.
func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	const q = `
		INSERT INTO users (
			user_id,
			username,
			is_active,
			team_name
		) VALUES ($1, $2, $3, $4)
	`

	exec := r.exec(ctx)

	_, err := exec.ExecContext(
		ctx,
		q,
		user.UserID,
		user.Username,
		user.IsActive,
		user.TeamName,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// unique_violation по PK user_id
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// UpdateSetIsActive — обновить is_active по user_id и вернуть обновлённого пользователя.
func (r *UserRepo) UpdateSetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	const q = `
		UPDATE users
		SET is_active = $2
		WHERE user_id = $1
		RETURNING user_id, username, is_active, team_name
	`

	exec := r.exec(ctx)

	var row userRow
	err := exec.
		QueryRowxContext(ctx, q, userID, isActive).
		StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// пользователя нет
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("update user is_active: %w", err)
	}

	return userFromRow(row), nil
}

// внутренний row-тип для users
type userRow struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
	TeamName string `db:"team_name"`
}

// маппинг row -> domain
func userFromRow(r userRow) domain.User {
	return domain.User{
		UserID:   r.UserID,
		Username: r.Username,
		IsActive: r.IsActive,
		TeamName: r.TeamName,
	}
}
