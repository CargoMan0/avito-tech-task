package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository"
	"github.com/CargoMan0/avito-tech-task/internal/service/dto"
	"github.com/google/uuid"
	"time"
)

func (s *Service) CreatePullRequest(ctx context.Context, data *dto.CreatePullRequestData) (*domain.PullRequest, error) {
	exists, err := s.pullRequestRepository.PullRequestExists(ctx, data.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("pull request repository: check if pull request exists: %w", err)
	}
	if exists {
		return nil, ErrPRAlreadyExists
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
			return nil, ErrNotFound
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
			return nil, uuid.Nil, ErrNotFound
		}

		return nil, uuid.Nil, fmt.Errorf("pull request repository: get pull request by reviewer id: %w", err)
	}

	seen := make(map[uuid.UUID]struct{})
	for _, reviewer := range pr.Reviewers {
		_, ok := seen[reviewer.ID]
		if ok {
			continue
		}

		seen[reviewer.ID] = struct{}{}
	}

	if _, ok := seen[oldReviewerID]; !ok {
		return nil, uuid.Nil, ErrUserNotAssignedToPR
	}

	if pr.Status == domain.PullRequestStatusMerged {
		return nil, uuid.Nil, ErrPRAlreadyMerged
	}

	user, err := s.userRepository.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, uuid.Nil, ErrNotFound
		}

		return nil, uuid.Nil, fmt.Errorf("user repository: get user by id: %w", err)
	}

	team, err := s.teamRepository.GetTeam(ctx, user.TeamName)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, uuid.Nil, ErrNotFound
		}

		return nil, uuid.Nil, fmt.Errorf("team repository: get team: %w", err)
	}

	// TODO: Fix wrong reviewer choosing here.

	suitableReviewers := make([]domain.User, 0)
	for _, teamUser := range team.Users {
		if !teamUser.IsActive {
			continue
		}
		if teamUser.ID == pr.AuthorID || teamUser.ID == oldReviewerID {
			continue
		}

		suitableReviewers = append(suitableReviewers, teamUser)
	}
	if len(suitableReviewers) == 0 {
		return nil, uuid.Nil, ErrPRNoSuitableCandidates
	}

	pr.Reviewers = suitableReviewers
	newReviewerID := chooseReviewersRandomly(suitableReviewers, 1)[0].ID

	err = s.pullRequestRepository.UpdatePullRequestReviewer(ctx, pr.ID, oldReviewerID, newReviewerID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("pull request repository: update pull request reviewer: %w", err)
	}

	return pr, newReviewerID, nil
}

func (s *Service) MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (*domain.PullRequest, error) {
	pr, err := s.pullRequestRepository.GetPullRequestByID(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, ErrNotFound
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
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("pull request repository: update pull request status: %w", err)
	}

	return pr, nil
}
