package main

import (
<<<<<<< HEAD
	"fmt"
	"peng-front/config"
	"peng-front/handlers/mattermost_command"
	"peng-front/handlers/rabbitmq_consumer"
	"peng-front/utils"
=======
    "log"
    "peng-front/config"
    "peng-front/handlers/rabbitmq_consumer"
    "peng-front/handlers/mattermost_command"
    "fmt"
>>>>>>> origin/main
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		utils.LogMessage(utils.ERROR, fmt.Sprintf("Failed to load configuration: %v", err))
	}

	if cfg.Mode == "Output" {
		r := rabbitmq_consumer.NewRabbitmqWorker(&cfg.RabbitMQ, "output", "output_mattermost")
		err = r.Start()
		if err != nil {
			utils.LogMessage(utils.ERROR, fmt.Sprintf("Failed to start worker: %v", err))
		}
		defer r.Stop()
		utils.LogMessage(utils.INFO, "Output Worker started")
	}

	if cfg.Mode == "Input" {
		m := mattermost_command.NewMattermostCommand(&cfg.Mattermost)
		m.Run(fmt.Sprintf(":%s", cfg.Server.Port))
		utils.LogMessage(utils.INFO, "Mattermost Command Handler started")
	}
}
