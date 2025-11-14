package handlers

import (
	"encoding/json"
	"errors"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"log"
	"net/http"
)

const (
	ErrCodeTeamExists       = "TEAM_EXISTS"
	ErrCodePRExists         = "PR_EXISTS"
	ErrCodePRMerged         = "PR_MERGED"
	ErrCodeNotAssigned      = "NOT_ASSIGNED"
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
)

func handleDomainError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, service.ErrTeamAlreadyExists):
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   ErrCodeTeamExists,
			Message: "team_name already exists",
		})
	case errors.Is(err, service.ErrPRAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   ErrCodePRExists,
			Message: "PR id already exists",
		})
	case errors.Is(err, service.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   ErrCodeResourceNotFound,
			Message: "resource not found",
		})
	case errors.Is(err, service.ErrUserNotAssignedToPR):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   ErrCodeNotAssigned,
			Message: "reviewer is not assigned to this PR",
		})
	case errors.Is(err, service.ErrPRAlreadyMerged):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   ErrCodePRMerged,
			Message: "cannot reassign on merged PR",
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func writeJSONError(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": errorMessage},
	)
}
