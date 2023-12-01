package main

import (
	"bytes"
	"encoding/json"
	"github.com/RakhimovAns/Time_Manager/types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	BotToken := "6791006120:AAFzk656CBCPbNlWVolFZl1cUwp4ej7A-Tc"
	//https://api.telegram.org/bot<token>/METHOD_NAME
	BotAPI := "https://api.telegram.org/bot"
	BotURL := BotAPI + BotToken
	offset := 0
	for {
		Updates, err := getUpdates(BotURL, offset)
		if err != nil {
			log.Println("Smth went wrong", err.Error())
		}
		for _, update := range Updates {
			err = respond(BotURL, update)
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(BotURL string, offset int) ([]types.Update, error) {
	resp, err := http.Get(BotURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponses types.RestResponse
	err = json.Unmarshal(body, &restResponses)
	if err != nil {
		return nil, err
	}
	return restResponses.Result, nil
}

func respond(botURL string, update types.Update) error {
	botMessage := new(types.BotMessage)
	botMessage.ChatId = update.Message.Chat.ChatId
	if update.Message.Text == "/help" {
		HelpRespond(botMessage)
	} else if update.Message.Text == "/author" {
		AuthorRespond(botMessage)
	} else if update.Message.Text == "/info" {
		InfoRespond(botMessage)
	} else {
		ErrorRespond(botMessage)
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func HelpRespond(botMessage *types.BotMessage) {
	botMessage.Text = "Hello, this bot can sort your doings and remind about them\n" +
		"You can use this following commands\n" +
		"/info - gets information about sorting methods\n" + //implemented
		"/sort - sorts your doings, example:\n" +
		"	/sort\n" +
		"	Make smth 01.12.2023 17:00\n" +
		"	Make smth2 02.12.2023 17:00\n" +
		"	...\n" +
		"/remind - reminds you about your doing, use this command like sort command\n" +
		"/author - gets information about authors\n" + //implemented
		"/delete - deletes doing from remind list,use this command like sort command\n" +
		"/list - gets all doing from remind list" // correct english grammar
	//"/change" //Add it later
}

func InfoRespond(botMessage *types.BotMessage) {
	botMessage.Text = "This bot sorts your doing by Eisenhower's Matrix.\n" + "Eisenhower's Matrix is the one of the most popular sorting methods of doing.The essence of the technique is to sort tasks by importance and urgency using a special table"
}

func AuthorRespond(botMessage *types.BotMessage) {
	botMessage.Text = "Ansar Rakhmimov. support: https://t.me/Rakhimov_Ans"
}

func ErrorRespond(botMessage *types.BotMessage) {
	botMessage.Text = "unrecognized command, use /help to get list of commands"
}
