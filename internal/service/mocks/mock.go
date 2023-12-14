package mocks

import (
	"context"

	"github.com/gogapopp/notificationService/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func NewMockRepositroy() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepository) InsertMessage(ctx context.Context, msg models.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockRepository) Subscribe(ctx context.Context, user models.UserSub) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Unsubscribe(ctx context.Context, userUnSub models.UserUnSub) error {
	args := m.Called(ctx, userUnSub)
	return args.Error(0)
}
