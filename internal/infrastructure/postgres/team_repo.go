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

// type TeamRepository interface {
// 	GetByUserID(ctx context.Context, userID string) (domain.Team, error)
// 	GetByTeamName(ctx context.Context, teamName string) (domain.Team, error)
// 	Create(ctx context.Context, team domain.Team) (domain.Team, error)
// }

type TeamRepo struct {
	db *sqlx.DB
}

func NewTeamRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) exec(ctx context.Context) dbExecutor {
	if tx := txmanager.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}

// Create — создать команду (без пользователей).
// Если команда с таким именем уже есть — domain.ErrAlreadyExists.
func (r *TeamRepo) Create(ctx context.Context, team domain.Team) (domain.Team, error) {
	const q = `
		INSERT INTO teams (team_name)
		VALUES ($1)
	`

	exec := r.exec(ctx)

	_, err := exec.ExecContext(ctx, q, team.TeamName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// unique_violation по PK team_name
			return domain.Team{}, domain.ErrAlreadyExists
		}

		return domain.Team{}, fmt.Errorf("insert team: %w", err)
	}

	// Возвращаем команду без членов — их создаёт отдельный user-репозиторий.
	return domain.Team{
		TeamName: team.TeamName,
		Members:  []domain.TeamMember{},
	}, nil
}

// GetByTeamName — получить команду и её участников по имени команды.
func (r *TeamRepo) GetByTeamName(ctx context.Context, teamName string) (domain.Team, error) {
	const q = `
	SELECT team_name, user_id, username, is_active
	FROM users
	WHERE team_name = $1
	`
	exec := r.exec(ctx)

	var mRows []teamMemberRow
	if err := exec.SelectContext(ctx, &mRows, q, teamName); err != nil {
		return domain.Team{}, fmt.Errorf("select team members: %w", err)
	}

	if len(mRows) == 0 {
		// нет ни одного участника с таким team_name
		return domain.Team{}, domain.ErrNotFound
	}

	return domain.Team{
		TeamName: teamName,
		Members:  membersFromRows(mRows),
	}, nil
}

// GetByUserID — найти команду по user_id и вернуть всю команду с её участниками.
func (r *TeamRepo) GetByUserID(ctx context.Context, userID string) (domain.Team, error) {
	const q = `
		SELECT team_name
		FROM users
		WHERE user_id = $1
	`

	exec := r.exec(ctx)

	var tr teamRow
	if err := exec.GetContext(ctx, &tr, q, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// пользователь не существует
			return domain.Team{}, domain.ErrNotFound
		}
		return domain.Team{}, fmt.Errorf("select team by user id: %w", err)
	}

	return r.GetByTeamName(ctx, tr.TeamName)
}

// row для команды (таблица teams)
type teamRow struct {
	TeamName string `db:"team_name"`
}

// row для участника команды (таблица users)
type teamMemberRow struct {
	TeamName string `db:"team_name"`
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
}

// маппинг members-строк к доменной модели
func membersFromRows(rows []teamMemberRow) []domain.TeamMember {
	if len(rows) == 0 {
		return []domain.TeamMember{}
	}

	res := make([]domain.TeamMember, 0, len(rows))
	for _, r := range rows {
		res = append(res, domain.TeamMember{
			IsActive: r.IsActive,
			UserID:   r.UserID,
			Username: r.Username,
		})
	}
	return res
}
