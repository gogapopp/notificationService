package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/gogapopp/notificationService/internal/delivery/httpserver"
	"github.com/gogapopp/notificationService/internal/logger"
	"github.com/gogapopp/notificationService/internal/repository/mongodb"
	"github.com/gogapopp/notificationService/internal/service"
	"github.com/labstack/echo"
)

func main() {
	ctx := context.Background()
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := mongodb.NewMongoDB(config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	service := service.NewService(db, logger)
	handler := httpserver.NewHandler(service, logger)

	e := echo.New()
	e.POST("/subscribe", handler.Subscribe)
	e.POST("/unsubscribe", handler.Unsubscribe)
	e.POST("/message", handler.Message)
	e.GET("/health", handler.Health)
	s := http.Server{
		Addr:         config.HTTPServer.Address,
		Handler:      e,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
