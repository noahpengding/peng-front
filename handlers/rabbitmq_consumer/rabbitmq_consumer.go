package rabbitmq_consumer

import (
    "encoding/json"
	"context"
    "peng-front/config"
	"peng-front/models"
	"peng-front/handlers/mattermost_poster"
    amqp "github.com/rabbitmq/amqp091-go"
    "fmt"
)

type RabbitmqWorker struct {
	conn	*amqp.Connection
	channel	*amqp.Channel
	query_name	string
	consumerTag	string
	topic          string
	groupID        string
	ctx            context.Context
	cancel         context.CancelFunc
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}

func NewRabbitmqWorker(config *config.RabbitMQConfig, topic string, groupID string) *RabbitmqWorker {
	conn, err := amqp.Dial(config.URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	err = ch.ExchangeDeclare(
		topic, // name
		"fanout", // type
		true, // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	query_name := fmt.Sprintf("%s.%s", topic, groupID)
	q, err := ch.QueueDeclare(
		query_name, // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-queue-type": "quorum",
			"x-consumer-group": groupID,
            "x-queue-leader-locator": "least-leaders",
		},
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name, // queue name
		"", // routing key
		topic, // exchange
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to bind a queue")
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &RabbitmqWorker{
		conn: conn,
		channel: ch,
		query_name: q.Name,
		consumerTag: "",
		topic: topic,
		groupID: groupID,
		ctx: ctx,
		cancel: cancel,
	}
}

func (w *RabbitmqWorker) Start() error {
	msgs, err := w.channel.Consume(
		w.query_name, // queue
		"", // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				fmt.Println(" [*] Worker stopped")
				return
			case msg, ok := <-msgs:
				if !ok {
					fmt.Println(" [*] msg is not ok")
					return
				}
				w.Handle_message(msg)
			}
		}
	}()
	fmt.Println(fmt.Sprintf(" [*] Waiting for messages in %s.%s|%s To exit press CTRL+C", w.topic, w.groupID, w.query_name))
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
		fmt.Println(err)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	mm := mattermost_poster.NewMattermostClient(&cfg.Mattermost)
	if data.Team == "" {
		data.Team = cfg.Mattermost.Team
	}
	if data.Channel == "" {
		data.Channel = cfg.Mattermost.Channel
	}
	err = mm.MattermostSend(data.Team, data.Channel, data.Data)
	if err != nil {
		fmt.Println(err)
	}
	msg.Ack(false)
}
