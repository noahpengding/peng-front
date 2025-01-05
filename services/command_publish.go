package services

import (
	"peng-front/config"
	"peng-front/handlers/rabbitmq_publisher"
	"peng-front/models"
	"peng-front/utils"

	"github.com/google/uuid"

	"fmt"
	"strings"
)

func CommandPublish(request *models.CommandRequest) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	command := strings.Split(request.Text, " ")
	message := &models.Message{}
	topic := ""
	if command[0] == "output" {
		topic = "output"
		message = &models.Message{
			ID:      uuid.New().String(),
			Topic:   topic,
			Data:    outputCommand(command[1:]),
			Channel: request.Channel,
			Team:    request.Team,
		}
	}

	if topic == "" || message == nil {
		topic = "output"
		message = &models.Message{
			ID:      uuid.New().String(),
			Topic:   topic,
			Data:    fmt.Sprintf("Invalid command: %s", request.Text),
			Channel: cfg.Mattermost.Channel,
			Team:    cfg.Mattermost.Team,
		}
	}

	r := rabbitmq_publisher.NewRabbitMQClient(&cfg.RabbitMQ)
	defer r.Close()
	if err := r.PublishMessage(topic, message); err != nil {
		utils.LogMessage(utils.WARN, fmt.Sprintf("Failed to publish message: %v", err))
	}
	return nil
}

func outputCommand(command []string) string {
	return strings.Join(command, " ")
}
