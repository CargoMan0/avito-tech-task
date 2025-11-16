package service

import (
	"github.com/CargoMan0/avito-tech-task/internal/repository"
)

type Service struct {
	pullRequestRepository repository.PullRequestRepository
	userRepository        repository.UserRepository
	teamRepository        repository.TeamRepository
}

func NewService(
	pullRequestsRepository repository.PullRequestRepository,
	usersRepository repository.UserRepository,
	teamsRepository repository.TeamRepository,
) *Service {
	return &Service{
		pullRequestRepository: pullRequestsRepository,
		userRepository:        usersRepository,
		teamRepository:        teamsRepository,
	}
}
