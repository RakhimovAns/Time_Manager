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

var English map[string]string

func SetVocabulary() {
	English["HelpRespond"] = "Hello, this bot can sort your doings and remind about them\\n\" +\n\t\t\"You can use this following commands\\n\" +\n\t\t\"/info - gets information about sorting methods\\n\" +\n\t\t\"/sort - sorts your doings, use this command in following format:\\n\" +\n\t\t\"\tName Date Time Importance(from 1 to 4, from lower to higher)\\n\" +\n\t\t\"\tExample: Task 8.02.2024 13:50 1\\n\" +\n\t\t\"/remind - reminds you about your doing, use this command like  a sort command\\n\" +\n\t\t\"/author - gets information about authors\\n\" +\n\t\t\"/delete - deletes doings from remind list,use this command like a sort command\\n\" +\n\t\t\"/list - gets all doings from remind list\\n\" +\n\t\t\"/done - you can use this command when you finished some doings, use this command like a sort command"
}
