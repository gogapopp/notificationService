package service

import (
	"context"
	"encoding/json"
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
	Unsubscribe(ctx context.Context, userUnSub models.UserUnSub) error
}

type Publisher interface {
	PublishMessage(body []byte) error
}

type Service struct {
	repo      Repository
	validate  *validator.Validate
	publisher Publisher
	log       *zap.SugaredLogger
}

func NewService(repository Repository, publisher Publisher, log *zap.SugaredLogger) *Service {
	return &Service{
		repo:      repository,
		validate:  validator.New(),
		publisher: publisher,
		log:       log,
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
	err = s.repo.InsertMessage(ctx, msg)
	if err != nil {
		return err
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = s.publisher.PublishMessage(msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Subscribe(ctx context.Context, userSub models.UserSub) error {
	err := s.validate.Struct(userSub)
	if err != nil {
		return ErrEmptyField
	}
	return s.repo.Subscribe(ctx, userSub)
}

func (s *Service) Unsubscribe(ctx context.Context, userUnSub models.UserUnSub) error {
	err := s.validate.Struct(userUnSub)
	if err != nil {
		return ErrEmptyField
	}
	return s.repo.Unsubscribe(ctx, userUnSub)
}
