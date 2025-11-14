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

type PullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{
		db: db,
	}
}

func (p *PullRequestRepository) PullRequestExists(ctx context.Context, pullRequestID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS (SELECT FROM pull_requests WHERE id = $1)`

	var exists bool
	err := p.db.QueryRowContext(ctx, query, pullRequestID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exec pull request exists sql query: %w", err)
	}

	return exists, nil
}

func (p *PullRequestRepository) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (err error) {
	const (
		insertPRQuery          = `INSERT INTO pull_requests (id, name, author_id, status, need_more_reviewers)  VALUES ($1, $2, $3, $4, $5);`
		insertPRReviewersQuery = `INSERT INTO pull_requests_to_reviewers (pull_request_id, user_id) VALUES %s;`
	)

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil && !errors.Is(err, sql.ErrTxDone) {
				err = errors.Join(err, fmt.Errorf("rollback: %w", rollbackErr))
			}
		}
	}()

	statusSQL := statusFromDomainToEnum(pr.Status)
	err = tx.QueryRowContext(ctx, insertPRQuery, pr.ID, pr.Name, pr.AuthorID, statusSQL, pr.NeedMoreReviewers).Err()
	if err != nil {
		return fmt.Errorf("run insert pull request sql query: %w", err)
	}

	if len(pr.Reviewers) > 0 {
		var values []interface{}
		var placeholders []string

		for i, reviewer := range pr.Reviewers {
			idx := i*2 + 1
			placeholders = append(placeholders, fmt.Sprintf("($%d,$%d)", idx, idx+1))
			values = append(values, pr.ID, reviewer.ID)
		}

		insertUserQuery := fmt.Sprintf(
			insertPRReviewersQuery,
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

func (p *PullRequestRepository) GetPullRequestByID(ctx context.Context, pullRequestID uuid.UUID) (*domain.PullRequest, error) {
	const query = `SELECT name, author_id, status, need_more_reviewers FROM pull_requests WHERE id = $1`

	pr := domain.PullRequest{
		ID: pullRequestID,
	}

	err := p.db.QueryRowContext(ctx, query, pullRequestID).Scan(&pr.Name, &pr.AuthorID, &pr.Status, &pr.NeedMoreReviewers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRepoNotFound
		}

		return nil, fmt.Errorf("run get pull request by id sql query: %w", err)
	}

	return &pr, nil
}

func (p *PullRequestRepository) GetPullRequestsByReviewerID(ctx context.Context, authorID uuid.UUID) (_ []domain.PullRequest, err error) {
	const query = `SELECT id, name, status, need_more_reviewers FROM pull_requests WHERE author_id = $1`

	res := make([]domain.PullRequest, 0)
	rows, err := p.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("run get pull requests by reviewers sql query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		pr := domain.PullRequest{}
		sc := PRStatusScanner{}

		err = rows.Scan(&pr.ID, &pr.Name, &sc, &pr.NeedMoreReviewers)
		if err != nil {
			return nil, fmt.Errorf("scan row into struct: %w", err)
		}

		pr.Status = sc.Status
		res = append(res, pr)
	}

	return res, nil
}

func (p *PullRequestRepository) UpdatePullRequestStatus(ctx context.Context, status domain.PullRequestStatus, pullRequestID uuid.UUID) error {
	const query = `UPDATE pull_requests SET status = $1 WHERE id = $2`

	statusSQL := statusFromDomainToEnum(status)

	_, err := p.db.ExecContext(ctx, query, statusSQL, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrRepoNotFound
		}

		return fmt.Errorf("exec update pull request status query: %w", err)
	}

	return nil
}

func (p *PullRequestRepository) UpdatePullRequestReviewer(ctx context.Context, pullRequestID uuid.UUID, reviewerID uuid.UUID) error {
	const query = `UPDATE pull_requests_to_reviewers SET user_id = $1 WHERE pull_request_id= $2`

	res, err := p.db.ExecContext(ctx, query, pullRequestID, reviewerID)
	if err != nil {
		return fmt.Errorf("exec update pull request reviewer sql query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrRepoNotFound
	}

	return nil
}
