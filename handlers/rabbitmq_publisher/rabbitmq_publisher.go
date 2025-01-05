package rabbitmq_publisher

import (
    "encoding/json"
    "peng-front/config"
	"peng-front/models"
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/google/uuid"

    "fmt"
	"context"
)

type RabbitmqClient struct {
    conn    *amqp.Connection
    channel *amqp.Channel  
    done    chan bool
}

func fail_on_error(err error, msg string) {
    if err != nil {
        fmt.Printf("%s: %s", msg, err)
    }
}

func NewRabbitMQClient(config *config.RabbitMQConfig) *RabbitmqClient {
    conn, err := amqp.Dial(config.URL)
    fail_on_error(err, "Failed to connect to RabbitMQ")
    ch, err := conn.Channel()
    fail_on_error(err, "Failed to open a channel")
    
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
        topic, // name
        "fanout", // type
        true, // durable
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil, // arguments
    )

    ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return c.channel.PublishWithContext(
		ctx,
		topic, // exchange
		"", // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			ContentType: "application/json",
			Body:        data,
		},
	)
}

func (c *RabbitmqClient) Close() {
	c.conn.Close()
	c.channel.Close()
}