package routes

import (
	"github.com/CargoMan0/avito-tech-task/internal/http/handlers"
	"github.com/CargoMan0/avito-tech-task/internal/http/middleware"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, handlers *handlers.Handlers) {
	mux.HandleFunc("POST /team/add", handlers.PostTeam())
	mux.HandleFunc("GET /team/get", middleware.AuthMiddleware(handlers.GetTeam()))

	mux.HandleFunc("POST /pullRequest/create", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostPullRequest())))
	mux.HandleFunc("POST /pullRequest/merge", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostPullRequestMerge())))
	mux.HandleFunc("POST /pullRequest/reassign", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostPullRequestReassign())))

	mux.HandleFunc("POST /users/setIsActive", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostUsersSetIsActive())))
	mux.HandleFunc("GET /users/getReview", middleware.AuthMiddleware(handlers.GetUsersReview()))
}
