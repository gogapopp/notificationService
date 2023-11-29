package service

import (
	"context"

	"github.com/gogapopp/notificationService/internal/repository/mongodb"
	"go.uber.org/zap"
)

type Repository interface {
	Ping(ctx context.Context) error
}

type Service struct {
	repo Repository
	log  *zap.SugaredLogger
}

func NewService(repository *mongodb.DB, log *zap.SugaredLogger) *Service {
	return &Service{repo: repository, log: log}
}

func (s *Service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
