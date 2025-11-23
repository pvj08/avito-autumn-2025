package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
)

// type PullRequestRepository interface {
// 	Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
// 	GetByID(ctx context.Context, id string) (domain.PullRequest, error)
// 	GetByUserID(ctx context.Context, userID string) ([]domain.PullRequest, error)
// 	Update(ctx context.Context, pr domain.PullRequest) error
// 	Save(ctx context.Context, pr domain.PullRequest) error
// }

type PullRequestRepo struct {
	db *sqlx.DB
}

func NewPullRequestRepo(db *sqlx.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) exec(ctx context.Context) dbExecutor {
	if tx := txmanager.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}

// Create — вставка нового PR с возвратом актуальных данных из БД.
func (r *PullRequestRepo) Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	const q = `
		INSERT INTO pull_requests (
			pull_request_id,
			author_id,
			pull_request_name,
			status,
			assigned_reviewers,
			created_at,
			merged_at
		) VALUES ($1, $2, $3, $4, $5, COALESCE($6, now()), $7)
		RETURNING pull_request_id,
		          author_id,
		          pull_request_name,
		          status,
		          assigned_reviewers,
		          created_at,
		          merged_at
	`

	row := toRow(pr)
	exec := r.exec(ctx)

	var dbRow pullRequestRow
	err := exec.
		QueryRowxContext(
			ctx,
			q,
			row.PullRequestID,
			row.AuthorID,
			row.PullRequestName,
			row.Status,
			row.AssignedReviewers,
			sql.NullTime{Time: row.CreatedAt, Valid: !row.CreatedAt.IsZero()},
			row.MergedAt,
		).
		StructScan(&dbRow)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 — уникальный ключ нарушен (duplicate key)
			if pgErr.Code == "23505" {
				return domain.PullRequest{}, domain.ErrAlreadyExists
			}
		}
		return domain.PullRequest{}, fmt.Errorf("insert pull_request: %w", err)
	}

	return fromRow(dbRow), nil
}

// GetByID — получить PR по ID.
func (r *PullRequestRepo) GetByID(ctx context.Context, id string) (domain.PullRequest, error) {
	const q = `
		SELECT pull_request_id,
		       author_id,
		       pull_request_name,
		       status,
		       assigned_reviewers,
		       created_at,
		       merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	exec := r.exec(ctx)

	var dbRow pullRequestRow
	err := exec.GetContext(ctx, &dbRow, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// если есть domain.ErrNotFound — можешь вернуть его
			return domain.PullRequest{}, domain.ErrNotFound
		}
		return domain.PullRequest{}, fmt.Errorf("select pull_request by id: %w", err)
	}

	return fromRow(dbRow), nil
}

// GetByUserID — PR'ы, назначенные конкретному ревьюверу (по assigned_reviewers).
func (r *PullRequestRepo) GetByUserID(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	const q = `
		SELECT pull_request_id,
		       author_id,
		       pull_request_name,
		       status,
		       assigned_reviewers,
		       created_at,
		       merged_at
		FROM pull_requests
		WHERE $1 = ANY(assigned_reviewers)
	`

	exec := r.exec(ctx)

	var rows []pullRequestRow
	if err := exec.SelectContext(ctx, &rows, q, userID); err != nil {
		return nil, fmt.Errorf("select pull_requests by user id: %w", err)
	}

	// нет строк — вернуть пустой массив, не nil
	if len(rows) == 0 {
		return []domain.PullRequest{}, nil
	}

	res := make([]domain.PullRequest, 0, len(rows))
	for _, row := range rows {
		res = append(res, fromRow(row))
	}

	return res, nil
}

// // Update — обновить существующий PR по его ID.
// func (r *PullRequestRepo) Update(ctx context.Context, pr domain.PullRequest) error {
// 	const q = `
// 		UPDATE pull_requests
// 		SET author_id          = $2,
// 		    pull_request_name  = $3,
// 		    status             = $4,
// 		    assigned_reviewers = $5,
// 		    created_at         = $6,
// 		    merged_at          = $7
// 		WHERE pull_request_id  = $1
// 	`

// 	row := toRow(pr)
// 	exec := r.exec(ctx)

// 	_, err := exec.ExecContext(
// 		ctx,
// 		q,
// 		row.PullRequestID,
// 		row.AuthorID,
// 		row.PullRequestName,
// 		row.Status,
// 		row.AssignedReviewers,
// 		row.CreatedAt,
// 		row.MergedAt,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("update pull_request: %w", err)
// 	}

// 	return nil
// }

// Save — upsert по pull_request_id.
func (r *PullRequestRepo) Save(ctx context.Context, pr domain.PullRequest) error {
	const q = `
		INSERT INTO pull_requests (
			pull_request_id,
			author_id,
			pull_request_name,
			status,
			assigned_reviewers,
			created_at,
			merged_at
		) VALUES ($1, $2, $3, $4, $5, COALESCE($6, now()), $7)
	`
	row := toRow(pr)
	exec := r.exec(ctx)

	_, err := exec.ExecContext(
		ctx,
		q,
		row.PullRequestID,
		row.AuthorID,
		row.PullRequestName,
		row.Status,
		row.AssignedReviewers,
		row.CreatedAt,
		row.MergedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// дубликат по primary key (pull_request_id)
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert pull_request: %w", err)
	}

	return nil
}

// внутренняя структура для маппинга в БД
type pullRequestRow struct {
	PullRequestID     string     `db:"pull_request_id"`
	AuthorID          string     `db:"author_id"`
	PullRequestName   string     `db:"pull_request_name"`
	Status            string     `db:"status"`
	AssignedReviewers []string   `db:"assigned_reviewers"`
	CreatedAt         time.Time  `db:"created_at"`
	MergedAt          *time.Time `db:"merged_at"`
}

// domain -> row
func toRow(pr domain.PullRequest) pullRequestRow {
	var createdAt time.Time
	if pr.CreatedAt != nil {
		createdAt = *pr.CreatedAt
	}

	return pullRequestRow{
		PullRequestID:     pr.PullRequestID,
		AuthorID:          pr.AuthorID,
		PullRequestName:   pr.PullRequestName,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          pr.MergedAt,
	}
}

// row -> domain
func fromRow(r pullRequestRow) domain.PullRequest {
	createdAt := r.CreatedAt

	return domain.PullRequest{
		PullRequestID:     r.PullRequestID,
		AuthorID:          r.AuthorID,
		PullRequestName:   r.PullRequestName,
		Status:            domain.PullRequestStatus(r.Status),
		AssignedReviewers: r.AssignedReviewers,
		CreatedAt:         &createdAt,
		MergedAt:          r.MergedAt,
	}
}
