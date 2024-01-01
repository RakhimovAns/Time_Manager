package main

import (
	"github.com/RakhimovAns/Time_Manager/pkg/Responds"
	"github.com/RakhimovAns/Time_Manager/pkg/postgresql"
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
	postgresql.ConnectToDB(DSN)
	pool := postgresql.GetPool()
	defer pool.Close()
	for {
		respond.StatusRespond(BotURL)
		respond.Remind(BotURL)
		respond.Want(BotURL)
		Updates, err := respond.GetUpdates(BotURL, offset)
		if err != nil {
			log.Println("Something went wrong", err.Error())
		}
		for _, update := range Updates {
			err = respond.Respond(BotURL, update)
			offset = update.UpdateId + 1
		}
	}
}
