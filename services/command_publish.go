package services

import (
	"peng-front/config"
	"peng-front/handlers/rabbitmq_publisher"
	"peng-front/models"
	"peng-front/utils"

	"github.com/google/uuid"

	"encoding/json"
	"fmt"
	"strings"
)

type Homelabcommand struct {
	Type    string      `json:"type" binding:"required"`
	Message interface{} `json:"message" binding:"required"`
	App     interface{} `json:"app"`
}

func CommandPublish(request *models.CommandRequest) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	command := strings.Split(request.Text, " ")
	message := &models.Message{}
	if command[0] == "output" {
		message = &models.Message{
			ID:      uuid.New().String(),
			Topic:   "output",
			Data:    outputCommand(command[1:]),
			Channel: request.Channel,
			Team:    request.Team,
		}
	}

	if command[0] == "homelab" {
		message = &models.Message{
			ID:      uuid.New().String(),
			Topic:   "homelab",
			Data:    homelabCommand(command[1:]),
			Channel: request.Channel,
			Team:    request.Team,
		}
	}

	if message == nil || message.Data == "" {
		return fmt.Errorf("Invalid command: %s", request.Text)
	}

	r := rabbitmq_publisher.NewRabbitMQClient(&cfg.RabbitMQ)
	defer r.Close()
	if err := r.PublishMessage(message.Topic, message); err != nil {
		utils.LogMessage(utils.WARN, fmt.Sprintf("Failed to publish message: %v", err))
	}
	return nil
}

func outputCommand(command []string) string {
	return strings.Join(command, " ")
}

func homelabCommand(command []string) *Homelabcommand {
	var message interface{}
	msgStr := strings.Join(command[1:], " ")
	if err := json.Unmarshal([]byte(msgStr), &message); err != nil {
		message = msgStr
	}
	return &Homelabcommand{
		Type:    command[0],
		Message: message,
	}
}
