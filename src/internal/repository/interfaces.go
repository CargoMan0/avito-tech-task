package repository

import (
	"context"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, data *domain.Team) error
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
	TeamExists(ctx context.Context, name string) (bool, error)
}
