package handlers

import (
	"encoding/json"
	"errors"
	"github.com/CargoMan0/avito-tech-task/src/internal/service"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, service.ErrTeamAlreadyExists):
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   "TEAM_EXISTS",
			"message": "team_name already exists",
		})
	case errors.Is(err, service.ErrPRAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   "PR_EXISTS",
			"message": "PR id already exists",
		})
	case errors.Is(err, service.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   "NOT_FOUND",
			"message": "resource not found",
		})
	}
}
