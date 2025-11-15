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

func pullRequestFromDomain(pullRequest *domain.PullRequest, mergedAt *time.Time) pullRequestDTO {
	pr := pullRequestDTO{
		PullRequestID:   pullRequest.ID.String(),
		PullRequestName: pullRequest.Name,
		Status:          convertPRStatusFromDomain(pullRequest.Status),
		AuthorID:        pullRequest.AuthorID.String(),
	}

	if mergedAt != nil {
		formatted := mergedAt.UTC().Format(time.RFC3339)
		pr.MergedAt = &formatted
	}

	reviewers := make([]string, 0, len(pullRequest.Reviewers))
	for _, reviewer := range pullRequest.Reviewers {
		reviewers = append(reviewers, reviewer.ID.String())
	}
	pr.AssignedReviewers = reviewers

	return pr
}

func tryConvertStatusToDomain(status string) (domain.PullRequestStatus, error) {
	switch status {
	case "MERGED":
		return domain.PullRequestStatusMerged, nil
	case "OPEN":
		return domain.PullRequestStatusOpen, nil
	default:
		return 0, fmt.Errorf("invalid string: (%s)", status)
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
