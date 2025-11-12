package handlers

import (
	"context"
	"encoding/json"
	serviceDTO "github.com/CargoMan0/avito-tech-task/src/internal/service/dto"
	"net/http"
)

type Service interface {
	GetTeam(ctx context.Context, name string) (*serviceDTO.TeamData, error)
	CreateTeam(ctx context.Context, data *serviceDTO.TeamData) error
}

func CreateTeam(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req teamDTO

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		err = validateTeamDTO(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := mapTeamDTOToService(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.CreateTeam(r.Context(), &data)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(req)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func GetTeam(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
