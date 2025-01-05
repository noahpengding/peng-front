package main

import (
    "log"
    "peng-front/config"
    "peng-front/handlers/rabbitmq_consumer"
    "peng-front/handlers/mattermost_command"

    "os"
    "fmt"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    if cfg.Mode == "Output" {
        r := rabbitmq_consumer.NewRabbitmqWorker(&cfg.RabbitMQ, "output", "output_mattermost")
        err = r.Start()
        if err != nil {
            fmt.Println(err)
        }
        defer r.Stop()

        fmt.Println("Worker started")
    }
    
    if cfg.Mode == "Input" {
        m := mattermost_command.NewMattermostCommand(&cfg.Mattermost)
        m.Run(fmt.Sprintf(":%s", cfg.Server.Port))
    }
}