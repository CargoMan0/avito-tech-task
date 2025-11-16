package handlers

import (
	"encoding/json"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"net/http"
)

func PostTeam(service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req teamDTO

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		err = validateTeamDTO(req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		team, err := teamDTOToDomain(req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := r.Context()
		err = service.CreateTeam(ctx, team)
		if err != nil {
			handleDomainError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(req)
		if err != nil {
			return
		}
	}
}

func GetTeam(service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")

		if len(teamName) == 0 || len(teamName) >= 20 {
			writeJSONError(w, http.StatusBadRequest, "invalid team_name")
			return
		}

		ctx := r.Context()
		team, err := service.GetTeam(ctx, teamName)
		if err != nil {
			handleDomainError(w, err)
			return
		}

		resp := teamDTO{
			Name:    team.Name,
			Members: make([]teamMemberDTO, 0, len(team.Users)),
		}
		for _, user := range team.Users {
			resp.Members = append(resp.Members, teamMemberDTO{
				UserID:   user.ID.String(),
				Username: user.Name,
				IsActive: user.IsActive,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}
