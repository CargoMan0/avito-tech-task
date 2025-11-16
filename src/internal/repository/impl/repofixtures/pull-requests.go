package repofixtures

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository/impl/repoconvert"
	"strings"
	"time"
)

type PullRequestFixture struct {
	db *sql.DB
}

func NewPullRequestFixture(db *sql.DB) *PullRequestFixture {
	return &PullRequestFixture{db}
}

func (p *PullRequestFixture) MustPrepareTestPullRequest(pr *domain.PullRequest) {
	const (
		insertPRQuery          = `INSERT INTO pull_requests (id, created_at, name, author_id, status, need_more_reviewers) VALUES ($1, $2, $3, $4, $5, $6);`
		insertPRReviewersQuery = `INSERT INTO pull_requests_to_reviewers (pull_request_id, user_id) VALUES %s;`
	)

	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("begin tx: %v", err))
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				panic(fmt.Sprintf("rollback failed: %v, original err: %v", rbErr, err))
			}
		}
	}()

	statusSQL := repoconvert.StatusFromDomainToEnum(pr.Status)
	if _, err = tx.ExecContext(ctx, insertPRQuery, pr.ID, time.Now(), pr.Name, pr.AuthorID, statusSQL, pr.NeedMoreReviewers); err != nil {
		panic(fmt.Sprintf("insert pull request failed: %v", err))
	}

	if len(pr.Reviewers) > 0 {
		var values []interface{}
		var placeholders []string

		for i, reviewer := range pr.Reviewers {
			idx := i*2 + 1
			placeholders = append(placeholders, fmt.Sprintf("($%d,$%d)", idx, idx+1))
			values = append(values, pr.ID, reviewer.ID)
		}

		insertUserQuery := fmt.Sprintf(insertPRReviewersQuery, strings.Join(placeholders, ","))
		if _, err = tx.ExecContext(ctx, insertUserQuery, values...); err != nil {
			panic(fmt.Sprintf("insert reviewers failed: %v", err))
		}
	}

	if err = tx.Commit(); err != nil {
		panic(fmt.Sprintf("commit tx failed: %v", err))
	}
}

func (p *PullRequestFixture) MustClearTestData() {
	const query = `
DELETE FROM pull_requests_to_reviewers;
DELETE FROM pull_requests;
`

	_, err := p.db.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("failed to clear pull request test data: %v", err))
	}
}
