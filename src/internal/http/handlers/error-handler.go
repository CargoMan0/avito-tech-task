package handlers

import (
	"encoding/json"
	"errors"
	httperrors "github.com/CargoMan0/avito-tech-task/internal/http/errors"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"log/slog"
	"net/http"
)

func handleDomainError(w http.ResponseWriter, err error, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, service.ErrTeamAlreadyExists):
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(httperrors.ErrorResponse{
			Error:   httperrors.ErrCodeTeamExists,
			Message: "team_name already exists",
		})
	case errors.Is(err, service.ErrPRAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(httperrors.ErrorResponse{
			Error:   httperrors.ErrCodePRExists,
			Message: "PR id already exists",
		})
	case errors.Is(err, service.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(httperrors.ErrorResponse{
			Error:   httperrors.ErrCodeResourceNotFound,
			Message: "resource not found",
		})
	case errors.Is(err, service.ErrUserNotAssignedToPR):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(httperrors.ErrorResponse{
			Error:   httperrors.ErrCodeNotAssigned,
			Message: "reviewer is not assigned to this PR",
		})
	case errors.Is(err, service.ErrPRAlreadyMerged):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(httperrors.ErrorResponse{
			Error:   httperrors.ErrCodePRMerged,
			Message: "cannot reassign on merged PR",
		})
	case errors.Is(err, service.ErrPRNoSuitableCandidates):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(httperrors.ErrorResponse{
			Error:   httperrors.ErrCodeNoCandidate,
			Message: "no active replacement candidate in team",
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("unexpected internal error from service",
			slog.String("error", err.Error()),
		)
	}
}
