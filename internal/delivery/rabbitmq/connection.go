package rabbitmq

import (
	"fmt"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/streadway/amqp"
)

func NewRabbitMQ(config *config.Config) (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQ.User, config.RabbitMQ.Password, config.RabbitMQ.Host, config.RabbitMQ.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = channel.QueueDeclare(
		config.RabbitMQ.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	return conn, err
}
