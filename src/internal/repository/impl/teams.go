package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository"
	"github.com/google/uuid"
	"strings"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (t *TeamRepository) CreateTeam(ctx context.Context, data *domain.Team) (err error) {
	const insertTeamQuery = `INSERT INTO teams (id, name) VALUES($1,$2)`

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = errors.Join(err, fmt.Errorf("rollback tx: %w", rollbackErr))
			}
		}
	}()

	teamID := uuid.New()
	_, err = tx.ExecContext(ctx, insertTeamQuery, teamID, data.Name)
	if err != nil {
		return fmt.Errorf("exec insert team sql query: %w", err)
	}

	// Batch insert team members. It was done, because there is no limit for users in team.
	if len(data.Users) > 0 {
		var values []interface{}
		var placeholders []string

		for i, user := range data.Users {
			idx := i*4 + 1
			placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", idx, idx+1, idx+2, idx+3))
			values = append(values, user.ID, user.Name, user.IsActive, teamID)
		}

		insertUserQuery := fmt.Sprintf(
			`INSERT INTO users (id,name,is_active,team_id) VALUES %s 
		            ON CONFLICT (id) DO UPDATE SET team_id = EXCLUDED.team_id;
                   `,
			strings.Join(placeholders, ","),
		)

		_, err = tx.ExecContext(ctx, insertUserQuery, values...)
		if err != nil {
			return fmt.Errorf("exec insert users batch sql query: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (t *TeamRepository) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	const (
		getTeamQuery  = `SELECT id FROM teams WHERE name = $1`
		getUsersQuery = `SELECT id, name, is_active FROM users WHERE team_id = $1`
	)

	team := domain.Team{
		Name: name,
	}

	var teamID uuid.UUID
	err := t.db.QueryRowContext(ctx, getTeamQuery, name).Scan(&teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRepoNotFound
		}

		return nil, fmt.Errorf("get team: %w", err)
	}

	rows, err := t.db.QueryContext(ctx, getUsersQuery, teamID)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.IsActive); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		team.Users = append(team.Users, u)
	}

	return &team, nil
}

func (t *TeamRepository) TeamExists(ctx context.Context, name string) (bool, error) {
	const query = `SELECT EXISTS(SELECT FROM teams WHERE name = $1)`

	var exists bool
	err := t.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("run sql query: %w", err)
	}

	return exists, nil
}
