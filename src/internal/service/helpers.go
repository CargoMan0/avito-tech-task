package service

import (
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"math/rand"
	"time"
)

func chooseReviewersRandomly(users []domain.User, maxReviewers int) []domain.User {
	candidates := make([]domain.User, 0, len(users))
	for _, u := range users {
		candidates = append(candidates, u)

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
