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

func (r *PullRequestRepo) Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	const q = `
		INSERT INTO pull_requests (pull_request_id, author_id, name, status, created_at, ...)
		VALUES (:pull_request_id, :author_id, :name, :status, NOW(), ...)
		RETURNING pull_request_id, author_id, name, status, created_at, ...;
	`

	row := fromDomain(pr)

	var dbPr pullReqRow
	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return domain.PullRequest{}, err
	}
	if err := stmt.GetContext(ctx, &dbPr, row); err != nil {
		// здесь разбираем ошибку Postgres
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			// 23505 — unique_violation
			return domain.PullRequest{}, ErrPullRequestAlreadyExists
		}
		return domain.PullRequest{}, err
	}

	return toDomain(dbPr), nil
}
