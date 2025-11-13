package routes

import (
	"github.com/CargoMan0/avito-tech-task/internal/http/handlers"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, svc *service.Service) {
	mux.HandleFunc("POST /team/add", handlers.PostTeam(svc))
	mux.HandleFunc("GET /team/get", handlers.GetTeam(svc))

	mux.HandleFunc("POST /pullRequest/create", handlers.PostPullRequest(svc))
	mux.HandleFunc("POST /pullRequest/merge", handlers.PostPullRequestMerge(svc))
	mux.HandleFunc("POST /pullRequest/reassign", handlers.PostPullRequestReassign(svc))

	mux.HandleFunc("POST /users/setIsActive", handlers.PostUsersSetIsActive(svc))
	mux.HandleFunc("GET /users/getReview", handlers.GetUsersReview(svc))
}
