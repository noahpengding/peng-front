package mattermost_poster

import (
	"peng-front/config"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

type MattermostClient struct {
	client *model.Client4
	config *config.MattermostConfig
}

func NewMattermostClient(config *config.MattermostConfig) *MattermostClient {
	client := model.NewAPIv4Client(config.URL)
	client.SetToken(config.Token)

	return &MattermostClient{
		client: client,
		config: config,
	}
}

func (c *MattermostClient) sendMessage(channelID string, message string) error {
	post := &model.Post{
		ChannelId: channelID,
		Message:   process_message(message),
	}

	if _, _, err := c.client.CreatePost(post); err != nil {
		return err
	}

	return nil
}

func process_message(message string) string {
	message = strings.ReplaceAll(message, "\\$", "")
	message = strings.ReplaceAll(message, "$", "")
	re := regexp.MustCompile(`\\\[\s*\n*`)
	message = re.ReplaceAllString(message, "$")
	re = regexp.MustCompile(`\n*\s*\\\]`)
	message = re.ReplaceAllString(message, "$")
	return message
}

func (c *MattermostClient) GetUserID(username string) (string, error) {
	user, _, err := c.client.GetUserByUsername(username, "")
	if err != nil {
		return "", err
	}
	return user.Id, nil
}

func (c *MattermostClient) getTeamID(teamName string) (string, error) {
	team, _, err := c.client.GetTeamByName(teamName, "")
	if err != nil {
		return "", err
	}
	return team.Id, nil
}

func (c *MattermostClient) getChannelID(channelName string, teamName string) (string, error) {
	channel, _, err := c.client.GetChannelByName(channelName, teamName, "")
	if err != nil {
		return "", err
	}
	return channel.Id, nil
}

func (c *MattermostClient) MattermostSend(team string, channel string, message string) error {
	teamID, err := c.getTeamID(team)
	if err != nil {
		return err
	}

	channelID, err := c.getChannelID(channel, teamID)
	if err != nil {
		return err
	}

	return c.sendMessage(channelID, message)
}

func (c *MattermostClient) MattermostSendWithChannelID(channelID string, message string) error {
	return c.sendMessage(channelID, message)
}
