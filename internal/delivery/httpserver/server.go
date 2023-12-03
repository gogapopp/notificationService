package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/gogapopp/notificationService/internal/models"
	"github.com/gogapopp/notificationService/internal/repository"
	"github.com/gogapopp/notificationService/internal/service"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Service interface {
	Ping(ctx context.Context) error
	InsertMessage(ctx context.Context, msg models.Message) error
	Subscribe(ctx context.Context, user models.UserSub) error
	Unsubscribe(ctx context.Context, userUnSub models.UserUnSub) error
}

type Handler struct {
	service Service
	log     *zap.SugaredLogger
}

func NewHandler(service Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{service: service, log: logger}
}

func (h *Handler) Health(c echo.Context) error {
	ctx := c.Request().Context()
	err := h.service.Ping(ctx)
	if err != nil {
		h.log.Error(err)
		return c.String(http.StatusServiceUnavailable, "fail!")
	}
	return c.String(http.StatusOK, "health: success!")
}

func (h *Handler) Message(c echo.Context) error {
	ctx := c.Request().Context()
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		h.log.Error(err)
		return c.String(http.StatusBadRequest, "invalid input")
	}
	err := h.service.InsertMessage(ctx, msg)
	if err != nil {
		h.log.Error(err)
		if errors.Is(err, service.ErrEmptyField) {
			return c.String(http.StatusBadRequest, "invalid input")
		}
		return c.String(http.StatusInternalServerError, "something went wrong")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "message: success",
	})
}

func (h *Handler) Subscribe(c echo.Context) error {
	ctx := c.Request().Context()
	var userSub models.UserSub
	if err := c.Bind(&userSub); err != nil {
		h.log.Error(err)
		return c.String(http.StatusBadRequest, "invalid input")
	}
	err := h.service.Subscribe(ctx, userSub)
	if err != nil {
		h.log.Error(err)
		// validate error
		if errors.Is(err, service.ErrEmptyField) {
			return c.String(http.StatusBadRequest, "invalid input")
		}
		if errors.Is(err, repository.ErrUserAlreadySubscribed) {
			return c.String(http.StatusBadRequest, "user already subscribed")
		}
		return c.String(http.StatusInternalServerError, "something went wrong")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "subscribe: success",
	})
}

func (h *Handler) Unsubscribe(c echo.Context) error {
	ctx := c.Request().Context()
	var userUnSub models.UserUnSub
	if err := c.Bind(&userUnSub); err != nil {
		h.log.Error(err)
		return c.String(http.StatusBadRequest, "invalid input")
	}
	err := h.service.Unsubscribe(ctx, userUnSub)
	if err != nil {
		h.log.Error(err)
		// validate error
		if errors.Is(err, service.ErrEmptyField) {
			return c.String(http.StatusBadRequest, "invalid input")
		}
		if errors.Is(err, repository.ErrUserNotExists) {
			return c.String(http.StatusBadRequest, "user doesn't exists")
		}
		return c.String(http.StatusInternalServerError, "something went wrong")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "unsubscribe: success",
	})
}
