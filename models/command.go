package models

type CommandRequest struct {
	Team    string `json:"team_domain"`
	Channel string `json:"channel_name"`
	User  string `json:"user_name"`
	Command   string `json:"command"`
	Text      string `json:"text"`
}

type CommandResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}