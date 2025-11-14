package handlers

import (
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/google/uuid"
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
