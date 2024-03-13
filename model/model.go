package model

import "time"

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Chat struct {
	ChatId int64 `json:"id"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Doing struct {
	ID         int
	ChatId     int64
	Name       string
	Data       time.Time
	Importance int
	Status     bool
}

//var Russian map[string]string
//
//func SetVocabulary() {
//
//}
