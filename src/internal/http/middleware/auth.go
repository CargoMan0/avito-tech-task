package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	httperrors "github.com/CargoMan0/avito-tech-task/internal/http/errors"
	"github.com/golang-jwt/jwt/v4"
	"log/slog"
	"net/http"
	"strings"
)

type contextKey string

const (
	CtxRoleKey contextKey = "role"

	RoleAdmin = "admin"
	RoleUser  = "user"
)

func AuthMiddleware(logger *slog.Logger, jwtSecret string) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				httperrors.WriteJSONError(w, http.StatusUnauthorized, "missing token")
				return
			}

			parts := strings.Split(tokenHeader, " ")
			if len(parts) != 2 {
				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid token format")
				return
			}
			if parts[0] != "Bearer" {
				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid token type")
				return
			}

			tokenStr := parts[1]

			parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err != nil {
				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			if !parsedToken.Valid {
				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				logger.Error("invalid token claims in auth middleware")

				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid claims")
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				logger.Error("failed to extract role from token in auth middleware")

				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			if role != RoleAdmin && role != RoleUser {
				logger.Error("unexpected role in auth middleware", slog.String("role", role))
				httperrors.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), CtxRoleKey, role)
			handlerFunc(w, r.WithContext(ctx))
		}
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
