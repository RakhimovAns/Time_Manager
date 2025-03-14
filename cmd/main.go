package main

import (
	"github.com/RakhimovAns/Time_Manager/pkg/postgresql"
	respond "github.com/RakhimovAns/Time_Manager/pkg/responds"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

const (
	BotToken = "6791006120:AAFzk656CBCPbNlWVolFZl1cUwp4ej7A-Tc"
	BotAPI   = "https://api.telegram.org/bot"
	BotURL   = BotAPI + BotToken
	DSN      = "postgresql://postgres:postgres@localhost:5433/postgres"
)

func main() {
	postgresql.ConnectToDB(DSN)
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}
	offset := 0
	log.Printf("Authorized on account %s", bot.Self.UserName)
	for {
		go func() {
			for {
				respond.StatusRespond(BotURL)
			}
		}()
		go func() {
			for {
				respond.Remind(BotURL)
			}
		}()
		go func() {
			for {
				respond.Want(BotURL)
			}
		}()
		u := tgbotapi.NewUpdate(offset)
		updates := bot.GetUpdatesChan(u)
		for update := range updates {
			if update.CallbackQuery != nil {
				callback := update.CallbackQuery
				if callback.Data == "ENG" {
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
					log.Println("зашел в eng")
					respond.SetLanguage("English", &msg)
					msg.Text = "Set English language"
					bot.Send(msg)
				} else if callback.Data == "RUS" {
					log.Println("into rus")
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
					respond.SetLanguage("Russian", &msg)
					msg.Text = "Поставлен Русский язык"
					bot.Send(msg)
				}
			}
			if update.Message != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				if update.Message.IsCommand() {
					if update.Message.Command() == "help" {
						respond.HelpRespond(&msg)
					} else if update.Message.Command() == "start" {
						respond.StartRespond(&msg)
					} else if update.Message.Command() == "info" {
						respond.InfoRespond(&msg)
					} else if update.Message.Command() == "author" {
						respond.AuthorRespond(&msg)
					} else if update.Message.Command() == "list" {
						respond.ListRespond(&msg)
					}
					bot.Send(msg)
				}
				if strings.HasPrefix(update.Message.Text, "/sort") {
					names := strings.TrimPrefix(update.Message.Text, "/sort")
					names = strings.TrimSpace(names)
					data := strings.Split(strings.Replace(names, "\r\n", "\n", -1), "\n")
					respond.SortRespond(data, &msg)
					bot.Send(msg)
				} else if strings.HasPrefix(update.Message.Text, "/remind") {
					names := strings.TrimPrefix(update.Message.Text, "/remind")
					names = strings.TrimSpace(names)
					data := strings.Split(strings.Replace(names, "\r\n", "\n", -1), "\n")
					respond.RemindRespond(data, &msg)
					bot.Send(msg)
				} else if strings.HasPrefix(update.Message.Text, "/delete") {
					names := strings.TrimPrefix(update.Message.Text, "/delete")
					names = strings.TrimSpace(names)
					data := strings.Split(strings.Replace(names, "\r\n", "\n", -1), "\n")
					respond.DeleteRespond(data, &msg)
					bot.Send(msg)
				} else if strings.HasPrefix(update.Message.Text, "/done") {
					names := strings.TrimPrefix(update.Message.Text, "/done")
					names = strings.TrimSpace(names)
					data := strings.Split(strings.Replace(names, "\r\n", "\n", -1), "\n")
					respond.DoneRespond(data, &msg)
					bot.Send(msg)
				}
				offset = update.UpdateID + 1
			}
		}
	}
}
