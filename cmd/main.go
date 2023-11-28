package main

import (
	"log"

	"github.com/gogapopp/notificationService/internal/config"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
}
