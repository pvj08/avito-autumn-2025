package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

type TeamRepo struct {
	db *sqlx.DB
}

func NewTeamRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, u domain.User) error {
	_, err := r.db.NamedExecContext(ctx, `
        INSERT INTO users (id, name) VALUES (:id, :name)
    `, u)
	return err
}
