package repository

import (
	"database/sql"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type pullRequestRepo struct {
	db  *sql.DB
	log logger.Logger
}

func NewPullRequestRepo(db *sql.DB, log logger.Logger) *pullRequestRepo {
	return &pullRequestRepo{
		db:  db,
		log: log,
	}
}

type pullReqRow struct {
	// TODO: определить структуру строки таблицы в бд
	ID int64 `db:"id"`
}

func (uc *pullRequestRepo) Create() {
	// TODO: Реализовать методы UC
}

func toDomain(row pullReqRow) domain.PullRequest {
	// TODO: сделать маппинг строки из бд в домеин
	return domain.PullRequest{}
}

func toDomainSlice(rows []pullReqRow) []domain.PullRequest {
	pulls := make([]domain.PullRequest, 0, len(rows))
	for _, row := range rows {
		pulls = append(pulls, toDomain(row))
	}
	return pulls
}
