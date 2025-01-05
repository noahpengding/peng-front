package models

type Message struct {
    ID      string `json:"id"`
    Topic   string `json:"topic"`
    Data    string `json:"data"`
    Channel string `json:"channel"`
    Team    string `json:"team"`
}