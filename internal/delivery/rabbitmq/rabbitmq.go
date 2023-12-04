package rabbitmq

import (
	"github.com/gogapopp/notificationService/internal/config"
	"github.com/streadway/amqp"
)

type (
	Publisher struct {
		channel *amqp.Channel
		queue   string
	}
	Consumer struct {
		channel *amqp.Channel
		queue   string
	}
)

func NewPublisher(conn *amqp.Connection, config *config.Config) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{channel: ch, queue: config.RabbitMQ.Queue}, nil
}

func (p *Publisher) PublishMessage(body []byte) error {
	err := p.channel.Publish(
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	return err
}

func NewConsumer(conn *amqp.Connection, config *config.Config) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{channel: ch, queue: config.RabbitMQ.Queue}, nil
}

func (c *Consumer) ConsumeMessages(handler func([]byte)) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()
	return nil
}
