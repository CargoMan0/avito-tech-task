package handlers

import (
	"encoding/json"
	"github.com/CargoMan0/avito-tech-task/internal/service/dto"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handlers) PostPullRequest() http.HandlerFunc {
	type response struct {
		Pr pullRequestDTO `json:"pr"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req postPullRequestCreateDTO

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid json body")
		}

		pullRequestID, err := uuid.Parse(req.PullRequestID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid pull_request_id")
			return
		}
		authorID, err := uuid.Parse(req.AuthorID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid pull_request_author_id")
			return
		}
		if len(req.PullRequestName) == 0 || len(req.PullRequestName) > 20 {
			writeJSONError(w, http.StatusBadRequest, "invalid pull_request_name")
			return
		}

		ctx := r.Context()
		data := &dto.CreatePullRequestData{
			PullRequestID:   pullRequestID,
			PullRequestName: req.PullRequestName,
			AuthorID:        authorID,
		}

		pr, err := h.service.CreatePullRequest(ctx, data)
		if err != nil {
			handleDomainError(w, err, h.logger)
			return
		}

		reviewers := make([]string, 0, len(pr.Reviewers))
		for _, reviewer := range pr.Reviewers {
			reviewers = append(reviewers, reviewer.ID.String())
		}

		resp := response{
			Pr: pullRequestFromDomain(pr),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}

func (h *Handlers) PostPullRequestMerge() http.HandlerFunc {
	type request struct {
		PullRequestID string `json:"pull_request_id"`
	}
	type response struct {
		Pr pullRequestDTO `json:"pr"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		prID, err := uuid.Parse(req.PullRequestID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid pull_request_id")
			return
		}

		pr, err := h.service.MergePullRequest(r.Context(), prID)
		if err != nil {
			handleDomainError(w, err, h.logger)
			return
		}

		resp := response{
			Pr: pullRequestFromDomain(pr),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}

func (h *Handlers) PostPullRequestReassign() http.HandlerFunc {
	type request struct {
		PullRequestId string `json:"pull_request_id"`
		OldReviewerId string `json:"old_reviewer_id"`
	}
	type response struct {
		Pr         pullRequestDTO `json:"pr"`
		ReplacedBy string         `json:"replaced_by"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		prID, err := uuid.Parse(req.PullRequestId)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid pull_request_id")
			return
		}
		oldReviewerID, err := uuid.Parse(req.OldReviewerId)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid old_reviewer_id")
			return
		}

		ctx := r.Context()
		pr, newReviewerID, err := h.service.ReassignPullRequestReviewer(ctx, prID, oldReviewerID)
		if err != nil {
			handleDomainError(w, err, h.logger)
			return
		}

		resp := response{
			Pr:         pullRequestFromDomain(pr),
			ReplacedBy: newReviewerID.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}
