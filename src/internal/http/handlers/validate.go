package handlers

import (
	"errors"
	"fmt"
	serviceDTO "github.com/CargoMan0/avito-tech-task/src/internal/service/dto"
	"github.com/google/uuid"
)

func validateTeamDTO(dto teamDTO) error {
	if len(dto.Name) == 0 || len(dto.Name) >= 20 {
		return errors.New("invalid team name")
	}

	for _, member := range dto.Members {
		if len(member.Username) == 0 || len(member.Username) >= 20 {
			return errors.New("invalid username")
		}

		if len(member.UserID) == 0 {
			return errors.New("user id can not be empty")
		}
	}

	return nil
}

func mapTeamDTOToService(teamDTO teamDTO) (serviceDTO.TeamData, error) {
	teamData := serviceDTO.TeamData{
		Name:  teamDTO.Name,
		Users: make([]serviceDTO.UserInfo, 0, len(teamDTO.Members)),
	}

	for _, member := range teamDTO.Members {
		userID, err := uuid.Parse(member.UserID)
		if err != nil {
			return serviceDTO.TeamData{}, fmt.Errorf("failed to parse user id: %w", err)
		}

		userInfo := serviceDTO.UserInfo{
			UserID:   userID,
			Username: member.Username,
			IsActive: member.IsActive,
		}

		teamData.Users = append(teamData.Users, userInfo)
	}

	return teamData, nil
}
