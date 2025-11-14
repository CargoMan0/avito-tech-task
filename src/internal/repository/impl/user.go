package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `SELECT u.name,
                          u.is_active,
                          t.name
                          FROM users u
                          JOIN teams t ON u.team_id = t.id
                          WHERE u.id = $1;`

	user := &domain.User{
		ID: id,
	}
	err := u.db.QueryRowContext(ctx, query, id).Scan(&user.Name, &user.IsActive, &user.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRepoNotFound
		}

		return nil, fmt.Errorf("exec get user by id sql query: %w", err)
	}

	return user, nil
}

func (u *UserRepository) UpdateUserIsActive(ctx context.Context, isActive bool, userID uuid.UUID) error {
	const query = `UPDATE users SET is_active = $1 WHERE id = $2`

	res, err := u.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return fmt.Errorf("exec update user is active sql query: %w", err)
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
