package middleware

import (
	"context"
	"encoding/json"
	httperrors "github.com/CargoMan0/avito-tech-task/internal/http/errors"
	"net/http"
)

type contextKey string

const (
	CtxRoleKey contextKey = "role"

	AdminToken = "hardcoded-admin-token"
	UserToken  = "hardcoded-user-token"

	RoleAdmin = "admin"
	RoleUser  = "user"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		var role string
		switch token {
		case AdminToken:
			role = RoleAdmin
		case UserToken:
			role = RoleUser
		default:
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxRoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AdminOnly(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(CtxRoleKey)
		if role != RoleAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(httperrors.ErrorResponse{
				Error:   httperrors.ErrCodeResourceNotFound,
				Message: "resource not found",
			})
			return
		}
		next.ServeHTTP(w, r)
	}
}
