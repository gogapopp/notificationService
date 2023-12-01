package httpserver

import (
	"context"
	"net/http"

	"github.com/gogapopp/notificationService/internal/models"
	"github.com/gogapopp/notificationService/internal/service"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Service interface {
	Ping(ctx context.Context) error
	InsertMessage(ctx context.Context, msg models.Message) error
}

type Handler struct {
	service Service
	log     *zap.SugaredLogger
}

func NewHandler(service *service.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{service: service, log: logger}
}

func (h *Handler) Health(c echo.Context) error {
	h.log.Info(c.Path())
	ctx := c.Request().Context()
	err := h.service.Ping(ctx)
	if err != nil {
		h.log.Error(err)
		return c.String(http.StatusServiceUnavailable, "DB: fail!")
	}
	return c.String(http.StatusOK, "DB: pong!")
}

func (h *Handler) Message(c echo.Context) error {
	h.log.Info(c.Path())
	ctx := c.Request().Context()
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		h.log.Error(err)
		return c.String(http.StatusBadRequest, "invalid input")
	}
	err := h.service.InsertMessage(ctx, msg)
	if err != nil {
		h.log.Error(err)
		return c.String(http.StatusInternalServerError, "failed insert to the DB")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (h *Handler) Subscribe(c echo.Context) error {
	h.log.Info(c.Path())
	return nil
}

func (h *Handler) Unsubscribe(c echo.Context) error {
	h.log.Info(c.Path())
	return nil
}
