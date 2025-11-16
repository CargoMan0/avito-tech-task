package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handlers) PostTeam() http.HandlerFunc {
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
		err = h.service.CreateTeam(ctx, team)
		if err != nil {
			handleDomainError(w, err, h.logger)
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

func (h *Handlers) GetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")

		if len(teamName) == 0 || len(teamName) >= 20 {
			writeJSONError(w, http.StatusBadRequest, "invalid team_name")
			return
		}

		ctx := r.Context()
		team, err := h.service.GetTeam(ctx, teamName)
		if err != nil {
			handleDomainError(w, err, h.logger)
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

		return
	}
}
