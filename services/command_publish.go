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
			User:    request.User,
			Channel: request.Channel,
			Team:    request.Team,
		}
	}

	if command[0] == "homelab" {
		message = &models.Message{
			ID:      uuid.New().String(),
			Topic:   "homelab",
			Data:    homelabCommand(command[1:]),
			User:    request.User,
			Channel: request.Channel,
			Team:    request.Team,
		}
	}

	if command[0] == "chat" {
		message = &models.Message{
			ID:      uuid.New().String(),
			Topic:   "chat",
			Data:    chatCommand(command[1:]),
			User:    request.User,
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

func homelabCommand(command []string) *models.Homelabcommand {
	var message interface{}
	msgStr := strings.Join(command[1:], " ")
	if err := json.Unmarshal([]byte(msgStr), &message); err != nil {
		message = msgStr
	}
	return &models.Homelabcommand{
		Type:    command[0],
		Message: message,
	}
}

func chatCommand(command []string) *models.Chatcommand {
	type_of_command := "chat"
	if strings.HasPrefix(command[0], "get_") || strings.HasPrefix(command[0], "set_") ||
		strings.HasPrefix(command[0], "image") || strings.HasPrefix(command[0], "list_") ||
		strings.HasPrefix(command[0], "end") || strings.HasPrefix(command[0], "chat") ||
		strings.HasPrefix(command[0], "index_") {
		type_of_command = command[0]
		command = command[1:]
	}
	operator := "openai"
	if len(command) > 0 &&
		(strings.HasPrefix(command[0], "gemini") || strings.HasPrefix(command[0], "claude") ||
			strings.HasPrefix(command[0], "openai")) || strings.HasPrefix(command[0], "rag") {
		operator = command[0]
		command = command[1:]
	}
	file_path := ""
	if len(command) > 0 && strings.HasPrefix(command[0], "--file=") {
		file_path = strings.Split(command[0], "=")[1]
		command = command[1:]
	}
	return &models.Chatcommand{
		Sources:   "mattermost",
		Type:      type_of_command,
		Operator:  operator,
		File_path: file_path,
		Message:   strings.Join(command, " "),
	}
}
