package routes

import (
	"github.com/CargoMan0/avito-tech-task/internal/http/handlers"
	"github.com/CargoMan0/avito-tech-task/internal/http/middleware"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, svc *service.Service) {
	mux.HandleFunc("POST /team/add", handlers.PostTeam(svc))
	mux.HandleFunc("GET /team/get", middleware.AuthMiddleware(handlers.GetTeam(svc)))

	mux.HandleFunc("POST /pullRequest/create", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostPullRequest(svc))))
	mux.HandleFunc("POST /pullRequest/merge", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostPullRequestMerge(svc))))
	mux.HandleFunc("POST /pullRequest/reassign", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostPullRequestReassign(svc))))

	mux.HandleFunc("POST /users/setIsActive", middleware.AuthMiddleware(middleware.AdminOnly(handlers.PostUsersSetIsActive(svc))))
	mux.HandleFunc("GET /users/getReview", middleware.AuthMiddleware(handlers.GetUsersReview(svc)))
}
