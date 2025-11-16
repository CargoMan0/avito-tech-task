package handlers

import (
	"errors"
)

func validateTeamDTO(dto teamDTO) error {
	if len(dto.Name) == 0 || len(dto.Name) >= 20 {
		return errors.New("invalid team name")
	}

	if len(dto.Members) == 0 {
		return errors.New("team must have at least one member")
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
