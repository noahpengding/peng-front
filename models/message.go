package models

type Message struct {
	ID      string      `json:"ID"`
	Topic   string      `json:"Topic"`
	Data    interface{} `json:"Data"`
	User    string      `json:"User"`
	Channel string      `json:"Channel"`
	Team    string      `json:"Team"`
}
