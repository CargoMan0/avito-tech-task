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
			writeJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid user_id")
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
	}
}

func GetUsersReview(svc *service.Service) http.HandlerFunc {
	type response struct {
		UserID       string                `json:"user_id"`
		PullRequests []pullRequestShortDTO `json:"pull_requests"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid user_id")
			return
		}

		ctx := r.Context()
		prs, err := svc.GetUserReviews(ctx, userID)
		if err != nil {
			handleDomainError(w, err)
			return
		}

		resp := response{
			UserID:       userID.String(),
			PullRequests: make([]pullRequestShortDTO, 0, len(prs)),
		}
		for _, pr := range prs {
			resp.PullRequests = append(resp.PullRequests, pullRequestShortFromDomain(&pr))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}
