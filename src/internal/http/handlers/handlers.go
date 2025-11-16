package handlers

import (
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"log/slog"
)

type Handlers struct {
	logger  *slog.Logger
	service *service.Service
}

func New(service *service.Service, logger *slog.Logger) *Handlers {
	return &Handlers{
		logger:  logger,
		service: service,
	}
}
