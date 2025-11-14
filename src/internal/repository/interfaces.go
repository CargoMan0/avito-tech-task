package repository

import (
	"context"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/google/uuid"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, data *domain.Team) error
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
	TeamExists(ctx context.Context, name string) (bool, error)
	GetTeamMembers(ctx context.Context, name string) ([]domain.User, error)
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUserIsActive(ctx context.Context, isActive bool, userID uuid.UUID) error
}

type PullRequestRepository interface {
	CreatePullRequest(ctx context.Context, data *domain.PullRequest) error
	PullRequestExists(ctx context.Context, pullRequestID uuid.UUID) (bool, error)
	GetPullRequestsByReviewerID(ctx context.Context, reviewerID uuid.UUID) ([]domain.PullRequest, error)
	GetPullRequestByID(ctx context.Context, pullRequestID uuid.UUID) (*domain.PullRequest, error)
	UpdatePullRequestStatus(ctx context.Context, status domain.PullRequestStatus, id uuid.UUID) error
	UpdatePullRequestReviewer(ctx context.Context, pullRequestID uuid.UUID, oldReviewerID uuid.UUID) error
}
