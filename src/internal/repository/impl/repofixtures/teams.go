package repofixtures

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/google/uuid"
	"strings"
)

type TeamFixture struct {
	db *sql.DB
}

func NewTeamFixture(db *sql.DB) *TeamFixture {
	return &TeamFixture{db}
}

func (t *TeamFixture) MustPrepareTestTeam(ctx context.Context, team *domain.Team) {
	const insertTeamQuery = `INSERT INTO teams (id, name) VALUES($1,$2)`

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("cannot prepare test team: %w", err))
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil && !errors.Is(err, sql.ErrTxDone) {
				err = errors.Join(err, fmt.Errorf("rollback tx: %w", rollbackErr))
			}
		}
	}()

	teamID := uuid.New() // Technical ID inside the database
	_, err = tx.Exec(insertTeamQuery, teamID, team.Name)
	if err != nil {
		panic(fmt.Sprintf("cannot prepare test team: %w", err))
	}

	if len(team.Users) > 0 {
		var values []interface{}
		var placeholders []string

		for i, user := range team.Users {
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

		_, err = tx.Exec(insertUserQuery, values...)
		if err != nil {
			panic(fmt.Sprintf("exec insert users batch sql query: %w", err))
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Sprintf("commit tx: %w", err))
	}
}
