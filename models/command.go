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

type Homelabcommand struct {
	Type    string      `json:"type" binding:"required"`
	Message interface{} `json:"message" binding:"required"`
	App     interface{} `json:"app"`
}

type Chatcommand struct {
	Type      string `json:"type" binding:"required"`
	Operator  string `json:"operator" binding:"required"`
	File_path string `json:"file_path"`
	Message   string `json:"message"`
}
