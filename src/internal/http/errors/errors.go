package errors

import (
	"encoding/json"
	"net/http"
)

const (
	ErrCodeTeamExists       = "TEAM_EXISTS"
	ErrCodePRExists         = "PR_EXISTS"
	ErrCodePRMerged         = "PR_MERGED"
	ErrCodeNotAssigned      = "NOT_ASSIGNED"
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
	ErrCodeNoCandidate      = "NO_CANDIDATE"
)

func WriteJSONError(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": errorMessage,
	},
	)
}
