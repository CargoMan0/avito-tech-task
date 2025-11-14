package handlers

import (
	"encoding/json"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"github.com/google/uuid"
	"net/http"
)

func PostUsersSetIsActive(svc *service.Service) http.HandlerFunc {
	type request struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	type response struct {
		User userDTO `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json body")
			return
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid user_id")
			return
		}

		ctx := r.Context()
		user, err := svc.SetUserIsActive(ctx, req.IsActive, userID)
		if err != nil {
			handleDomainError(w, err)
			return
		}

		resp := response{
			User: userDTO{
				UserId:   user.UserID.String(),
				Username: user.Username,
				IsActive: req.IsActive,
				TeamName: user.TeamName,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}

		return
	}
}

func GetUsersReview(svc *service.Service) http.HandlerFunc {
	type response struct {
		UserId       string           `json:"user_id"`
		PullRequests []pullRequestDTO `json:"pull_requests"`
	}

	return func(w http.ResponseWriter, r *http.Request) {}
}
