package service

import "errors"

// Service Errors
var (
	ErrTeamAlreadyExists = errors.New("team already exists")

	ErrPRAlreadyExists        = errors.New("PR already exists")
	ErrPRAlreadyMerged        = errors.New("PR already merged")
	ErrPRNoSuitableCandidates = errors.New("PR no suitable candidates")

	ErrUserNotAssignedToPR = errors.New("user not assigned to PR")
	ErrUserDuplicated      = errors.New("user duplicated")

	ErrNotFound = errors.New("not found")
)
