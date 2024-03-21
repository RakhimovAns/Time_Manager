package parser

import (
	"errors"
	"github.com/RakhimovAns/Time_Manager/model"
	"github.com/RakhimovAns/Time_Manager/pkg/postgresql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

func Parse(data []string, botMessage *tgbotapi.MessageConfig) ([]model.Doing, error) {

	var Doings []model.Doing

	for _, doing := range data {

		SplitedData := strings.Fields(doing)
		if len(SplitedData) != 4 {
			language := postgresql.GetLanguage(botMessage.ChatID)
			botMessage.Text = "invalid type of doing, use command /help to get more information"
			if language == "Russian" {
				botMessage.Text = Translate(botMessage.Text)
			}
			return nil, errors.New("invalid type of doing")
		}

		var Do model.Doing

		Do.Name = SplitedData[0]
		DateTimeStr := SplitedData[1] + " " + SplitedData[2]
		layout := "2.01.2006 15:04"

		dateTime, err := time.Parse(layout, DateTimeStr)
		if err != nil {
			language := postgresql.GetLanguage(botMessage.ChatID)
			botMessage.Text = "invalid type of doing, use command /help to get more information"
			if language == "Russian" {
				botMessage.Text = Translate(botMessage.Text)
			}
			return nil, errors.New("invalid type of doing")
		}

		Do.Data = dateTime
		Do.Importance, _ = strconv.Atoi(SplitedData[3])
		Doings = append(Doings, Do)
	}
	return Doings, nil
}
