package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u domain.User) error {
	_, err := r.db.NamedExecContext(ctx, `
        INSERT INTO users (id, name) VALUES (:id, :name)
    `, u)
	return err
}
