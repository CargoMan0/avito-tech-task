package service

import "errors"

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrPRAlreadyExists   = errors.New("PR already exists")

	ErrNotFound = errors.New("not found")
)
