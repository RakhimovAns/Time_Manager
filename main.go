package main

import (
	"github.com/RakhimovAns/Time_Manager/pkg/Responds"
	"log"
)

const (
	BotToken = "6791006120:AAFzk656CBCPbNlWVolFZl1cUwp4ej7A-Tc"
	BotAPI   = "https://api.telegram.org/bot"
	BotURL   = BotAPI + BotToken
	DSN      = "postgresql://postgres:postgres@localhost:5432/manager"
)

func main() {
	offset := 0
	Responds.ConnectToDB(DSN)
	pool := Responds.GetPool()
	defer pool.Close()
	for {
		Responds.Remind(BotURL)
		Updates, err := Responds.GetUpdates(BotURL, offset)
		if err != nil {
			log.Println("Something went wrong", err.Error())
		}
		for _, update := range Updates {
			err = Responds.Respond(BotURL, update)
			offset = update.UpdateId + 1
		}
	}
}
