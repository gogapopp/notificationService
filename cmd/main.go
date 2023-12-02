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
	"github.com/labstack/echo/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger, err := logger.NewLogger()
	if err != nil {
		logger.Fatal(err)
	}
	config, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}
	db, err := mongodb.NewMongoDB(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := db.Client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	service := service.NewService(db, logger)
	handler := httpserver.NewHandler(service, logger)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/health", handler.Health)
	e.POST("/message", handler.Message)
	e.POST("/subscribe", handler.Subscribe)
	e.POST("/unsubscribe", handler.Unsubscribe)
	s := http.Server{
		Handler:      e,
		Addr:         config.HTTPServer.Address,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	logger.Infof("server is running at %s address", config.HTTPServer.Address)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
