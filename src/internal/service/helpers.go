package service

import (
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func chooseReviewers(users []domain.User, authorID uuid.UUID, maxReviewers int) []domain.User {
	candidates := make([]domain.User, 0, len(users))
	for _, u := range users {
		if u.ID != authorID || u.IsActive {
			candidates = append(candidates, u)
		}
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > maxReviewers {
		candidates = candidates[:maxReviewers]
	}

	return candidates
}
