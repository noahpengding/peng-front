package mattermost_command

import (
	"peng-front/config"
	"peng-front/models"
	"peng-front/services"
	"peng-front/utils"

	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"
)

func NewMattermostCommand(config *config.MattermostConfig) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/peng", handler_command)

	return router
}

func handler_command(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, models.CommandResponse{
			ResponseType: "ephemeral",
			Text:         "Invalid request",
		})
		return
	}
	utils.LogMessage(utils.DEBUG, fmt.Sprintf("Request body: %v", requestBody))

	var command models.CommandRequest
	if err := c.BindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, models.CommandResponse{
			ResponseType: "ephemeral",
			Text:         "Invalid request",
		})
		return
	}

	if !authentication(c.GetHeader("Authorization")) {
		c.JSON(http.StatusUnauthorized, models.CommandResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Unauthorized: %s", c.GetHeader("Authorization")),
		})
		return
	}

	if command.Command == "/peng" {
		response := process_command(command)
		c.JSON(http.StatusOK, response)
	}
}

func authentication(token string) bool {
	cfg, err := config.Load()
	if err != nil {
		return false
	}
	return fmt.Sprintf("Token %s", cfg.Mattermost.Command_Token) == token
}

func process_command(command models.CommandRequest) models.CommandResponse {
	if err := services.CommandPublish(&command); err != nil {
		return models.CommandResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Error: %s", err),
		}
	}
	return models.CommandResponse{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Received command: %s; Text: %s", command.Command, command.Text),
	}
}
