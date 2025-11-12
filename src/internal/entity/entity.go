package entity

import "github.com/google/uuid"

type Team struct {
	Name  string // Unique
	Users []User
}
type User struct {
	ID       uuid.UUID // Unique, immutable
	Name     string
	IsActive bool
}

type PullRequest struct {
	ID                uuid.UUID // Unique, immutable
	AuthorID          uuid.UUID
	Name              string
	Status            PullRequestStatus
	Reviewers         []User
	NeedMoreReviewers bool
}

type PullRequestStatus uint8

const (
	PullRequestStatusOpen PullRequestStatus = iota + 1
	PullRequestStatusMerged
)
