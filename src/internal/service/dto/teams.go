package dto

import "github.com/google/uuid"

type TeamData struct {
	Name  string
	Users []UserInfo
}

type UserInfo struct {
	UserID   uuid.UUID
	Username string
	IsActive bool
}
