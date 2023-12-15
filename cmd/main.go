package main

import (
	"context"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/gogapopp/notificationService/internal/delivery/httpserver"
	"github.com/gogapopp/notificationService/internal/delivery/rabbitmq"
	"github.com/gogapopp/notificationService/internal/logger"
	"github.com/gogapopp/notificationService/internal/repository/mongodb"
	"github.com/gogapopp/notificationService/internal/service"
	"github.com/gogapopp/notificationService/internal/service/email"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger, err := logger.NewLogger()
	fatal(err)
	config, err := config.NewConfig()
	fatal(err)
	db, err := mongodb.NewMongoDB(config)
	fatal(err)
	defer func() {
		if err := db.Client.Disconnect(ctx); err != nil {
			fatal(err)
		}
	}()
	rabbitmqConnection, err := rabbitmq.NewRabbitMQ(config)
	fatal(err)
	defer func() {
		if err := rabbitmqConnection.Close(); err != nil {
			fatal(err)
		}
	}()
	publisher, err := rabbitmq.NewPublisher(rabbitmqConnection, config)
	fatal(err)
	consumer, err := rabbitmq.NewConsumer(rabbitmqConnection, config)
	fatal(err)
	service := service.NewService(db, publisher, logger)
	emailAuth := smtp.PlainAuth("", config.SMTP.Login, config.SMTP.Password, config.SMTP.Provider)
	mailService := email.NewMailService(db, config, emailAuth, logger)
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
			fatal(err)
		}
	}()
	logger.Infof("server is running at %s address", config.HTTPServer.Address)
	go func() {
		if err := consumer.ConsumeMessages(mailService.SendEMails); err != nil {
			fatal(err)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	if err := s.Shutdown(ctx); err != nil {
		fatal(err)
	}
}
