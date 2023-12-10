package types

import "time"

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}
type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Doing struct {
	Name       string
	Data       time.Time
	Importance int
}
