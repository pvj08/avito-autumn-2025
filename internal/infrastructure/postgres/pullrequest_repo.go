package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

type PullRequestRepo struct {
	db *sqlx.DB
}

func NewPullRequestRepo(db *sqlx.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) Create(ctx context.Context, u domain.User) error {
	_, err := r.db.NamedExecContext(ctx, `
        INSERT INTO users (id, name) VALUES (:id, :name)
    `, u)
	return err
}
