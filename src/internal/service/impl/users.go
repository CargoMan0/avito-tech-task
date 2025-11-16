package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"github.com/CargoMan0/avito-tech-task/internal/service/impl/dto"
	"github.com/google/uuid"
)

func (s *Service) SetUserIsActive(ctx context.Context, isActive bool, userID uuid.UUID) (*dto.SetUserIsActiveResponse, error) {
	err := s.userRepository.UpdateUserIsActive(ctx, isActive, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("user repository: update user is active: %w", err)
	}

	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("user repository: get user: %w", err)
	}

	resp := dto.SetUserIsActiveResponse{
		UserID:   user.ID,
		Username: user.Name,
		IsActive: user.IsActive,
		TeamName: user.TeamName,
	}

	return &resp, nil
}

func (s *Service) GetUserReviews(ctx context.Context, userID uuid.UUID) ([]domain.PullRequest, error) {
	exists, err := s.userRepository.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user repository: user exists: %w", err)
	}
	if !exists {
		return nil, service.ErrNotFound
	}

	prs, err := s.pullRequestRepository.GetPullRequestsByReviewerID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("pull request repository: get pull requests: %w", err)
	}

	return prs, nil
}
