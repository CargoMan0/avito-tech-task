package routes

import (
	"github.com/CargoMan0/avito-tech-task/internal/http/handlers"
	"github.com/CargoMan0/avito-tech-task/internal/http/middleware"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, handlers *handlers.Handlers, authMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("POST /team/add", handlers.PostTeam())
	mux.HandleFunc("GET /team/get", authMiddleware(handlers.GetTeam()))

	mux.HandleFunc("POST /pullRequest/create", authMiddleware(middleware.AdminOnly(handlers.PostPullRequest())))
	mux.HandleFunc("POST /pullRequest/merge", authMiddleware(middleware.AdminOnly(handlers.PostPullRequestMerge())))
	mux.HandleFunc("POST /pullRequest/reassign", authMiddleware(middleware.AdminOnly(handlers.PostPullRequestReassign())))

	mux.HandleFunc("POST /users/setIsActive", authMiddleware(middleware.AdminOnly(handlers.PostUsersSetIsActive())))
	mux.HandleFunc("GET /users/getReview", authMiddleware(handlers.GetUsersReview()))
}
