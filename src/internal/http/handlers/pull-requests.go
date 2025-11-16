package handlers

import (
	"encoding/json"
	httperrors "github.com/CargoMan0/avito-tech-task/internal/http/errors"
	"github.com/CargoMan0/avito-tech-task/internal/service/impl/dto"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (h *Handlers) PostPullRequest() http.HandlerFunc {
	type request struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}
	type response struct {
		Pr pullRequestDTO `json:"pr"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid json body")
		}

		pullRequestID, err := uuid.Parse(req.PullRequestID)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid pull_request_id")
			return
		}
		authorID, err := uuid.Parse(req.AuthorID)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid pull_request_author_id")
			return
		}
		if len(req.PullRequestName) == 0 || len(req.PullRequestName) > 20 {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid pull_request_name")
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

		resp := response{
			Pr: pullRequestFromDomain(pr),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.logger.Error("failed to encode create pull request response",
				slog.String("error", err.Error()),
			)

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
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		prID, err := uuid.Parse(req.PullRequestID)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid pull_request_id")
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
			h.logger.Error("failed to encode merge pull request response",
				slog.String("error", err.Error()),
			)

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
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		prID, err := uuid.Parse(req.PullRequestId)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid pull_request_id")
			return
		}
		oldReviewerID, err := uuid.Parse(req.OldReviewerId)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusBadRequest, "invalid old_reviewer_id")
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
			h.logger.Error("failed to encode reassign pull request response",
				slog.String("error", err.Error()),
			)

			return
		}
	}
}
