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
	"time"
)

func (s *Service) CreatePullRequest(ctx context.Context, data *dto.CreatePullRequestData) (*domain.PullRequest, error) {
	exists, err := s.pullRequestRepository.PullRequestExists(ctx, data.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("pull request repository: check if pull request exists: %w", err)
	}
	if exists {
		return nil, service.ErrPRAlreadyExists
	}

	pr := &domain.PullRequest{
		ID:       data.PullRequestID,
		AuthorID: data.AuthorID,
		Name:     data.PullRequestName,
		Status:   domain.PullRequestStatusOpen,
	}

	user, err := s.userRepository.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, service.ErrNotFound
		}

		return nil, fmt.Errorf("user repository: get user: %w", err)
	}

	team, err := s.teamRepository.GetTeam(ctx, user.TeamName)
	if err != nil {
		return nil, fmt.Errorf("team repository: get team by user id: %w", err)
	}

	reviewers := make([]domain.User, 0)
	for _, teamUser := range team.Users {
		if teamUser.ID == data.AuthorID || !teamUser.IsActive {
			continue
		}

		reviewers = append(reviewers, teamUser)
	}

	pr.Reviewers = chooseReviewersRandomly(reviewers, domain.MaxReviewers)
	pr.NeedMoreReviewers = pr.CheckIfNeedMoreReviewers()

	err = s.pullRequestRepository.CreatePullRequest(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("pull request repository: create pull request: %w", err)
	}

	return pr, nil
}

func (s *Service) ReassignPullRequestReviewer(ctx context.Context, pullRequestID, oldReviewerID uuid.UUID) (*domain.PullRequest, uuid.UUID, error) {
	pr, err := s.pullRequestRepository.GetPullRequestByID(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, uuid.Nil, service.ErrNotFound
		}
		return nil, uuid.Nil, fmt.Errorf("pull request repository: get pull request: %w", err)
	}

	currentReviewers := make(map[uuid.UUID]struct{})
	for _, r := range pr.Reviewers {
		currentReviewers[r.ID] = struct{}{}
	}

	if _, ok := currentReviewers[oldReviewerID]; !ok {
		return nil, uuid.Nil, service.ErrUserNotAssignedToPR
	}

	if pr.Status == domain.PullRequestStatusMerged {
		return nil, uuid.Nil, service.ErrPRAlreadyMerged
	}

	author, err := s.userRepository.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, uuid.Nil, service.ErrNotFound
		}
		return nil, uuid.Nil, fmt.Errorf("user repository: get author: %w", err)
	}

	team, err := s.teamRepository.GetTeam(ctx, author.TeamName)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, uuid.Nil, service.ErrNotFound
		}
		return nil, uuid.Nil, fmt.Errorf("team repository: get team: %w", err)
	}

	var suitable []domain.User
	for _, u := range team.Users {
		if !u.IsActive || u.ID == pr.AuthorID || u.ID == oldReviewerID {
			continue
		}
		if _, exists := currentReviewers[u.ID]; exists {
			continue
		}
		suitable = append(suitable, u)
	}

	if len(suitable) == 0 {
		return nil, uuid.Nil, service.ErrPRNoSuitableCandidates
	}

	newReviewer := chooseReviewersRandomly(suitable, 1)[0]

	err = s.pullRequestRepository.UpdatePullRequestReviewer(ctx, pr.ID, oldReviewerID, newReviewer.ID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("pull request repository: update reviewer: %w", err)
	}

	for i, r := range pr.Reviewers {
		if r.ID == oldReviewerID {
			pr.Reviewers[i].ID = newReviewer.ID
			break
		}
	}

	return pr, newReviewer.ID, nil
}

func (s *Service) MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (*domain.PullRequest, error) {
	pr, err := s.pullRequestRepository.GetPullRequestByID(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("pull request repository: get pull request by id: %w", err)
	}

	if pr.Status == domain.PullRequestStatusMerged {
		return pr, nil
	}

	mergedAt := time.Now()
	pr.MergedAt = &mergedAt
	pr.Status = domain.PullRequestStatusMerged

	err = s.pullRequestRepository.UpdatePullRequestStatusAndMergedAt(ctx, domain.PullRequestStatusMerged, pullRequestID, mergedAt)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("pull request repository: update pull request status: %w", err)
	}

	return pr, nil
}
