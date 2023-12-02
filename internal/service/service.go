package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gogapopp/notificationService/internal/models"
	"go.uber.org/zap"
)

var ErrEmptyField = errors.New("empty field")

type Repository interface {
	Ping(ctx context.Context) error
	InsertMessage(ctx context.Context, msg models.Message) error
	Subscribe(ctx context.Context, user models.UserSub) error
}

type Service struct {
	repo     Repository
	validate *validator.Validate
	log      *zap.SugaredLogger
}

func NewService(repository Repository, log *zap.SugaredLogger) *Service {
	return &Service{
		repo:     repository,
		validate: validator.New(),
		log:      log,
	}
}

func (s *Service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *Service) InsertMessage(ctx context.Context, msg models.Message) error {
	err := s.validate.Struct(msg)
	if err != nil {
		return ErrEmptyField
	}
	msg.Timestamp = time.Now()
	return s.repo.InsertMessage(ctx, msg)
}

func (s *Service) Subscribe(ctx context.Context, userSub models.UserSub) error {
	err := s.validate.Struct(userSub)
	if err != nil {
		return ErrEmptyField
	}
	return s.repo.Subscribe(ctx, userSub)
}
