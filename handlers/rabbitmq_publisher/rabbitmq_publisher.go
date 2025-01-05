package rabbitmq_publisher

import (
	"encoding/json"
	"peng-front/config"
	"peng-front/models"
	"peng-front/utils"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"context"
	"fmt"
)

type RabbitmqClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan bool
}

func NewRabbitMQClient(config *config.RabbitMQConfig) *RabbitmqClient {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Connecting to RabbitMQ at %s with error %s", config.URL, err))
	}
	ch, err := conn.Channel()
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Opening channel with error %s", err))
	}
	return &RabbitmqClient{
		conn:    conn,
		channel: ch,
		done:    make(chan bool),
	}
}

func (c *RabbitmqClient) PublishMessage(topic string, message *models.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	c.channel.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return c.channel.PublishWithContext(
		ctx,
		topic, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			ContentType:  "application/json",
			Body:         data,
		},
	)
}

func (c *RabbitmqClient) Close() {
	c.conn.Close()
	c.channel.Close()
}
