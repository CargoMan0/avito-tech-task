package service

import (
	"context"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/service/impl/dto"
	"github.com/google/uuid"
)

// Service на текущий момент не требует разбиения на более маленькие и специфичные. Учитывая текущий размер системы и ТЗ,
// на текущий момент один интерфейс вполне подходит.
type Service interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, name string) (*domain.Team, error)

	SetUserIsActive(ctx context.Context, isActive bool, userID uuid.UUID) (*dto.SetUserIsActiveResponse, error)
	GetUserReviews(ctx context.Context, userID uuid.UUID) ([]domain.PullRequest, error)

	MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (*domain.PullRequest, error)
	CreatePullRequest(ctx context.Context, data *dto.CreatePullRequestData) (*domain.PullRequest, error)
	ReassignPullRequestReviewer(ctx context.Context, pullRequestID, oldReviewerID uuid.UUID) (*domain.PullRequest, uuid.UUID, error)
}
