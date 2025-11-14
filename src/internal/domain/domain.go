package domain

import "github.com/google/uuid"

type Team struct {
	Name  string // Unique
	Users []User // Team members
}

type User struct {
	ID       uuid.UUID // Unique, immutable
	TeamName string
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

const MaxReviewers = 2

func (pr *PullRequest) CheckIfNeedMoreReviewers() bool {
	return len(pr.Reviewers) < MaxReviewers
}

type PullRequestStatus uint8

const (
	PullRequestStatusOpen PullRequestStatus = iota + 1
	PullRequestStatusMerged
)
