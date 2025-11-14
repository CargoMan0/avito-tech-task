package domain

import "github.com/google/uuid"

type Team struct {
	Name  string // Unique
	Users []User
}

type User struct {
	ID       uuid.UUID // Unique, immutable
	TeamName string    // Name of the team, which user
	Name     string
	IsActive bool
}

type PullRequest struct {
	ID                uuid.UUID // Unique, immutable
	AuthorID          uuid.UUID // Immutable
	Name              string
	Status            PullRequestStatus
	Reviewers         []User
	NeedMoreReviewers bool
}

func NeedMoreReviewers(users []User) bool {
	return len(users) < 3
}

type PullRequestStatus uint8

const (
	PullRequestStatusOpen PullRequestStatus = iota + 1
	PullRequestStatusMerged
)
