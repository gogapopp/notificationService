package service

import (
	"context"
	"testing"
	"time"

	"github.com/gogapopp/notificationService/internal/models"
	"github.com/gogapopp/notificationService/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishMessage(msg []byte) error {
	args := m.Called(msg)
	return args.Error(0)
}

// success cases
func TestPing(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	service := NewService(repo, nil, nil)
	repo.On("Ping", mock.Anything).Return(nil)
	err := service.Ping(context.Background())
	repo.AssertCalled(t, "Ping", mock.Anything)
	assert.NoError(t, err)
}

func TestInsertMessage(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	publisher := new(MockPublisher)
	service := NewService(repo, publisher, nil)
	msg := models.Message{
		UserID:    "1",
		Message:   "test message",
		Timestamp: time.Now(),
	}
	publisher.On("PublishMessage", mock.Anything).Return(nil)
	repo.On("InsertMessage", mock.Anything, msg).Return(nil)
	err := service.InsertMessage(context.Background(), msg)
	repo.AssertCalled(t, "InsertMessage", mock.Anything, msg)
	assert.NoError(t, err)
}

func TestSubscribe(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	service := NewService(repo, nil, nil)
	userSub := models.UserSub{
		UserID: "1",
		Email:  "test@example.com",
	}
	repo.On("Subscribe", mock.Anything, userSub).Return(nil)
	err := service.Subscribe(context.Background(), userSub)
	repo.AssertCalled(t, "Subscribe", mock.Anything, userSub)
	assert.NoError(t, err)
}

func TestUnsubscribeSuccess(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	service := NewService(repo, nil, nil)
	userUnSub := models.UserUnSub{
		UserID: "1",
	}
	repo.On("Unsubscribe", mock.Anything, mock.Anything).Return(nil)
	err := service.Unsubscribe(context.Background(), userUnSub)
	repo.AssertCalled(t, "Unsubscribe", mock.Anything, mock.Anything)
	assert.NoError(t, err)
}

// error cases
func TestInsertMessageEmptyFields(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	service := NewService(repo, nil, nil)
	msg := models.Message{}
	err := service.InsertMessage(context.Background(), msg)
	repo.AssertNotCalled(t, "InsertMessage", mock.Anything, msg)
	assert.EqualError(t, err, "empty field")
}

func TestSubscribeEmptyFields(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	service := NewService(repo, nil, nil)
	userSub := models.UserSub{}
	err := service.Subscribe(context.Background(), userSub)
	repo.AssertNotCalled(t, "Subscribe", mock.Anything, userSub)
	assert.EqualError(t, err, "empty field")
}

func TestUnsubscribeEmptyFields(t *testing.T) {
	repo := mocks.NewMockRepositroy()
	service := NewService(repo, nil, nil)
	userUnSub := models.UserUnSub{}
	err := service.Unsubscribe(context.Background(), userUnSub)
	repo.AssertNotCalled(t, "Unsubscribe", mock.Anything, userUnSub)
	assert.EqualError(t, err, "empty field")
}
