package models

type CommandRequest struct {
	Team    string `form:"team_domain" json:"team_domain"`
	Channel string `form:"channel_name" json:"channel_name"`
	User    string `form:"user_name" json:"user_name"`
	Command string `form:"command" json:"command" binding:"required"`
	Text    string `form:"text" json:"text"`
}

type CommandResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}
