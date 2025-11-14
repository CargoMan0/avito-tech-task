package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository"
	"github.com/google/uuid"
)

func (s *Service) CreateTeam(ctx context.Context, team *domain.Team) error {
	seen := make(map[uuid.UUID]struct{})
	for _, user := range team.Users {
		if _, ok := seen[user.ID]; ok {
			return ErrUserDuplicated
		}

		seen[user.ID] = struct{}{}
	}

	exists, err := s.teamRepository.TeamExists(ctx, team.Name)
	if err != nil {
		return fmt.Errorf("team repository: team exists: %w", err)
	}
	if exists {
		return ErrTeamAlreadyExists
	}

	err = s.teamRepository.CreateTeam(ctx, team)
	if err != nil {
		return fmt.Errorf("team repository: create team: %w", err)
	}

	return nil
}

func (s *Service) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	team, err := s.teamRepository.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("team repository: get team: %w", err)
	}
	return team, nil
}
