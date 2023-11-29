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
	var botMessage types.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	if update.Message.Text == "/help" {
		botMessage.Text = "Use that commands"
	} else {
		botMessage.Text = "Smth went wrong"
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
