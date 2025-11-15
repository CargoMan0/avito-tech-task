package dto

import "github.com/google/uuid"

type CreatePullRequestData struct {
	PullRequestID   uuid.UUID
	AuthorID        uuid.UUID
	PullRequestName string
}

type SetUserIsActiveResponse struct {
	UserID   uuid.UUID
	Username string
	TeamName string
	IsActive bool
}
