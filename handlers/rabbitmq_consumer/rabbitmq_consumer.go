package rabbitmq_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"peng-front/config"
	"peng-front/handlers/mattermost_poster"
	"peng-front/models"
	"peng-front/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqWorker struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	consumerTag string
	topic       string
	groupID     string
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewRabbitmqWorker(config *config.RabbitMQConfig, topic string, groupID string) *RabbitmqWorker {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Connecting to RabbitMQ at %s with error %s", config.URL, err))
	}
	ch, err := conn.Channel()
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Opening channel with error %s", err))
	}
	err = ch.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Declared exchange %s with error %s", topic, err))
	}
	q, err := ch.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Declared queue %s with error %s", topic, err))
	}
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Bound queue %s to exchange %s with error %s", q.Name, topic, err))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &RabbitmqWorker{
		conn:        conn,
		channel:     ch,
		consumerTag: "",
		topic:       topic,
		groupID:     groupID,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (w *RabbitmqWorker) Start() error {
	msgs, err := w.channel.Consume(
		w.topic, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Failed to register a consumer: %s", err))
	} else {
		utils.LogMessage(utils.INFO, fmt.Sprintf(" [*] Waiting for messages in %s To exit press CTRL+C", w.topic))
	}

	var forever chan struct{}
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				utils.LogMessage(utils.INFO, fmt.Sprintf(" [*] Worker stopped in %s", w.topic))
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				w.Handle_message(msg)
			}
		}
	}()
	<-forever

	return nil
}

func (w *RabbitmqWorker) Stop() {
	w.cancel()
	w.conn.Close()
	w.channel.Close()
}

func (w *RabbitmqWorker) Handle_message(msg amqp.Delivery) {
	data := &models.Message{}
	err := json.Unmarshal(msg.Body, data)
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Failed to unmarshal message: %s", err))
	}
	cfg, err := config.Load()
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Failed to load config: %v", err))
	}
	mm := mattermost_poster.NewMattermostClient(&cfg.Mattermost)
	if data.Team == "" {
		data.Team = cfg.Mattermost.Team
	}
	if data.Channel == "" {
		data.Channel = cfg.Mattermost.Channel
	}
	utils.LogMessage(utils.DEBUG, fmt.Sprintf("Team: %s, Channel: %s, Data: %s", data.Team, data.Channel, data.Data))
	err = mm.MattermostSend(data.Team, data.Channel, data.Data.(string))
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Failed to send message: %s", err))
	} else {
		utils.LogMessage(utils.INFO, fmt.Sprintf("Sent message to %s in %s", data.Channel, data.Team))
	}
	msg.Ack(false)
}
