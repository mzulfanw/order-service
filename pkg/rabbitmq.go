package pkg

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mzulfanw/order-service/configs"
	"github.com/streadway/amqp"
)

type RabbitMQPublisher interface {
	PublishEvent(event string, payload interface{})
	Close()
}

type rabbitMQ struct {
	channel *amqp.Channel
}

func NewRabbitMQ(config configs.RabbitMQConfig) RabbitMQPublisher {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.User, config.Pass, config.Host, config.Port)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
	}

	ch, _ := conn.Channel()
	ch.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	return &rabbitMQ{channel: ch}
}

func (r rabbitMQ) PublishEvent(event string, payload interface{}) {
	body, _ := json.Marshal(payload)
	err := r.channel.Publish("events", event, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Printf("failed publish: %v", err)
	}
}

func (r rabbitMQ) Close() {
	r.channel.Close()
}
