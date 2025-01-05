package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Mode       string
	Server     ServerConfig
	Mattermost MattermostConfig
	RabbitMQ   RabbitMQConfig
}

type ServerConfig struct {
	Port     string
	GIN_MODE string
}

type MattermostConfig struct {
	URL           string
	Token         string
	Command_Token string
	Team          string
	Channel       string
}

type RabbitMQConfig struct {
	URL         string
	ExchangeOut string
}

func Load() (*Config, error) {
	config := &Config{}

	// Load from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigName("config_sample")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	if mode := os.Getenv("APP_MODE"); mode != "" {
		config.Mode = mode
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if gin_mode := os.Getenv("GIN_MODE"); gin_mode != "" {
		config.Server.GIN_MODE = gin_mode
	}
	if url := os.Getenv("MATTERMOST_URL"); url != "" {
		config.Mattermost.URL = url
	}
	if team := os.Getenv("MATTERMOST_TEAM"); team != "" {
		config.Mattermost.Team = team
	}
	if channel := os.Getenv("MATTERMOST_CHANNEL"); channel != "" {
		config.Mattermost.Channel = channel
	}
	if token := os.Getenv("MATTERMOST_TOKEN"); token != "" {
		config.Mattermost.Token = token
	}
	if command_token := os.Getenv("MATTERMOST_COMMAND_TOKEN"); command_token != "" {
		config.Mattermost.Command_Token = command_token
	}
	if url := os.Getenv("RABBITMQ_URL"); url != "" {
		config.RabbitMQ.URL = url
	}
	if exchangeOut := os.Getenv("RABBITMQ_EXCHANGE_OUT"); exchangeOut != "" {
		config.RabbitMQ.ExchangeOut = exchangeOut
	}

	return config, nil
}
