package handlers

import (
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/google/uuid"
	"time"
)

func teamDTOToDomain(teamDTO teamDTO) (*domain.Team, error) {
	team := domain.Team{
		Name:  teamDTO.Name,
		Users: make([]domain.User, 0, len(teamDTO.Members)),
	}

	for _, member := range teamDTO.Members {
		user, err := teamMemberDTOToDomain(member)
		if err != nil {
			return nil, err
		}
		team.Users = append(team.Users, user)
	}

	return &team, nil
}

func teamMemberDTOToDomain(dto teamMemberDTO) (domain.User, error) {
	userID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to parse user id: %w", err)
	}

	user := domain.User{
		ID:       userID,
		Name:     dto.Username,
		IsActive: dto.IsActive,
	}

	return user, nil
}

func pullRequestFromDomain(pullRequest *domain.PullRequest) pullRequestDTO {
	pr := pullRequestDTO{
		PullRequestID:   pullRequest.ID.String(),
		PullRequestName: pullRequest.Name,
		Status:          convertPRStatusFromDomain(pullRequest.Status),
		AuthorID:        pullRequest.AuthorID.String(),
	}

	if pullRequest.MergedAt != nil {
		pr.MergedAt = pullRequest.MergedAt.UTC().Format(time.RFC3339)
	}

	reviewers := make([]string, len(pullRequest.Reviewers))
	for i, reviewer := range pullRequest.Reviewers {
		reviewers[i] = reviewer.ID.String()
	}
	pr.AssignedReviewers = reviewers

	return pr
}

func pullRequestShortFromDomain(pullRequest *domain.PullRequest) pullRequestShortDTO {
	return pullRequestShortDTO{
		PullRequestID:   pullRequest.ID.String(),
		PullRequestName: pullRequest.Name,
		Status:          convertPRStatusFromDomain(pullRequest.Status),
		AuthorID:        pullRequest.AuthorID.String(),
	}
}

func convertPRStatusFromDomain(status domain.PullRequestStatus) string {
	switch status {
	case domain.PullRequestStatusMerged:
		return "MERGED"
	case domain.PullRequestStatusOpen:
		return "OPEN"
	default:
		panic(fmt.Sprintf("invalid domain.PullRequestStatus: (%v)", status))
	}
}
