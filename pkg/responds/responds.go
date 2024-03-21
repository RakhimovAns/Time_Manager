package respond

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RakhimovAns/Time_Manager/model"
	"github.com/RakhimovAns/Time_Manager/pkg/parser"
	"github.com/RakhimovAns/Time_Manager/pkg/postgresql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func SetLanguage(language string, config *tgbotapi.MessageConfig) {
	fmt.Println(language)
	if language != "English" && language != "Russian" {
		config.Text = "invalid language"
		return
	}
	postgresql.SetLanguage(config.ChatID, language)
	text := "Was set successfully"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	config.Text = text
}
func HelpRespond(config *tgbotapi.MessageConfig) {
	language := postgresql.GetLanguage(config.ChatID)
	text := "Hello, this bot can sort your doings and remind about them\n" + "You can use this following commands\n" + "/info - gets information about sorting methods\n" + "/sort - sorts your doings, use this command in following format:\n" + "Name Date Time Importance(from 1 to 4, from lower to higher)\n" + "Example: Task 8.02.2024 13:50 1\n" + "/remind - reminds you about your doing, use this command like  a sort command\n" + "/author - gets information about authors\n" + "/delete - deletes doings from remind list,use this command like a sort command\n" + "/list - gets all doings from remind list\n" + "/done - you can use this command when you finished some doings, use this command like a sort command\n"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	config.Text = text
}
func StartRespond(config *tgbotapi.MessageConfig) {
	err := postgresql.GetLanguageStatus(config.ChatID)
	if err != nil {
		inlineBtn := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("ENGLISH", "ENG"),
				tgbotapi.NewInlineKeyboardButtonData("Русский", "RUS"),
			),
		)
		config.ReplyMarkup = inlineBtn
		config.Text = ""
		return
	}
	language := postgresql.GetLanguage(config.ChatID)
	text := "Hello dear user! This bot sorts your doings, to get more info use command /help"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	config.Text = text
}
func InfoRespond(config *tgbotapi.MessageConfig) {
	language := postgresql.GetLanguage(config.ChatID)
	text := "This bot sorts your doing by ABCDE method.\n" + "ABCDE method is the one of the most popular sorting methods of doing.The essence of the technique is to sort tasks by importance  using a special table"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	config.Text = text
}

func DoneRespond(data []string, botMessage *tgbotapi.MessageConfig) {
	Doings, err := parser.Parse(data, botMessage)
	if err != nil {
		return
	}
	for _, doing := range Doings {
		postgresql.SetStatus(doing)
	}
	language := postgresql.GetLanguage(botMessage.ChatID)
	text := "Command finished successfully"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	botMessage.Text = text
}

func StatusRespond(botURL string) {
	Doings := postgresql.GetDoingsWithStatus()
	set := make(map[int64]time.Time)
	for _, doing := range Doings {
		set[doing.ChatId] = doing.Data
	}
	for chat_id, timer := range set {
		if time.Now().Add(3*time.Hour).Sub(timer).Hours() == 1 {
			var botMessage model.BotMessage
			botMessage.ChatId = chat_id
			language := postgresql.GetLanguage(chat_id)
			text := "Have you done anything from your doing list?"
			if language == "Russian" {
				text = parser.Translate(text)
			}
			botMessage.Text = text
			buf, err := json.Marshal(botMessage)
			if err != nil {
				log.Fatal(err)
			}
			_, err = http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
func AuthorRespond(config *tgbotapi.MessageConfig) {
	language := postgresql.GetLanguage(config.ChatID)
	text := "Ansar Rakhmimov. support: @Rakhimov_Ans"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	config.Text = text
}

func ErrorRespond(config *tgbotapi.MessageConfig) {
	config.Text = "unrecognized command, use /help to get list of commands"
	language := postgresql.GetLanguage(config.ChatID)
	if language == "Russian" {
		config.Text = parser.Translate(config.Text)
	}
}

func SortRespond(data []string, botMessage *tgbotapi.MessageConfig) {
	Doings, err := parser.Parse(data, botMessage)
	if err != nil {
		return
	}
	sort.SliceStable(Doings, func(i, j int) bool {
		if Doings[i].Data != Doings[j].Data {
			return Doings[i].Data.Before(Doings[j].Data)
		}
		return Doings[i].Importance > Doings[j].Importance
	})
	var answer string
	for _, item := range Doings {
		answer += item.Name + " " + item.Data.Format("2.01.2006 15:04") + " " + strconv.Itoa(item.Importance) + "\n"
	}
	botMessage.Text = answer
}

func RemindRespond(data []string, botMessage *tgbotapi.MessageConfig) {
	Doings, err := parser.Parse(data, botMessage)
	if err != nil {
		return
	}
	sort.SliceStable(Doings, func(i, j int) bool {
		if Doings[i].Data != Doings[j].Data {
			return Doings[i].Data.Before(Doings[j].Data)
		}
		return Doings[i].Importance > Doings[j].Importance
	})
	for _, doing := range Doings {
		err := postgresql.AddDoings(doing, botMessage)
		if err != nil {
			log.Fatal("problem with adding to db")
		}
		language := postgresql.GetLanguage(botMessage.ChatID)
		text := "I will remind about that"
		if language == "Russian" {
			text = parser.Translate(text)
		}
		botMessage.Text = text
	}
}

func Remind(botURL string) {
	Doings := postgresql.GetDoings()
	for _, doing := range Doings {
		var BotMessage model.BotMessage
		BotMessage.ChatId = doing.ChatId
		if doing.Data.Sub(time.Now().Add(3*time.Hour)) <= 0 {
			postgresql.SetStatus(doing)
			language := postgresql.GetLanguage(BotMessage.ChatId)
			text := "You need to start '"
			if language == "Russian" {
				text = parser.Translate(text)
			}
			BotMessage.Text = text + doing.Name + "'"
			buf, err := json.Marshal(BotMessage)
			if err != nil {
				log.Fatal(err)
			}
			_, err = http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ListRespond(config *tgbotapi.MessageConfig) {
	Doings := postgresql.GetDoingsByID(config.ChatID)
	if len(Doings) == 0 {
		language := postgresql.GetLanguage(config.ChatID)
		text := "You have no doings"
		if language == "Russian" {
			text = parser.Translate(text)
		}
		config.Text = text
		return
	}
	var answer string
	sort.SliceStable(Doings, func(i, j int) bool {
		if Doings[i].Data != Doings[j].Data {
			return Doings[i].Data.Before(Doings[j].Data)
		}
		return Doings[i].Importance > Doings[j].Importance
	})
	for _, item := range Doings {
		answer += item.Name + " " + item.Data.Format("2.01.2006 15:04") + " " + strconv.Itoa(item.Importance) + "\n"
	}
	config.Text = answer
}

func DeleteRespond(data []string, botMessage *tgbotapi.MessageConfig) {
	Doings, err := parser.Parse(data, botMessage)
	if err != nil {
		return
	}
	for _, doing := range Doings {
		err := postgresql.Delete(doing)
		if err != nil {
			log.Fatal("problem to delete")
		}
	}
	language := postgresql.GetLanguage(botMessage.ChatID)
	text := "Was deleted successfully"
	if language == "Russian" {
		text = parser.Translate(text)
	}
	botMessage.Text = text
}

func Want(botURL string) {
	if !(time.Now().Hour() == 22 && time.Now().Minute() == 30) {
		return
	}
	Doings := postgresql.GetAllDoings()
	set := make(map[int64]int)
	for _, doing := range Doings {
		if doing.Status == true {
			set[doing.ChatId]++
			postgresql.Delete(doing)
		}
	}
	type user struct {
		chat_id int64
		count   int
	}
	users := []user{}
	for chat_id, count := range set {
		users = append(users, user{chat_id: chat_id, count: count})
		botMessage := model.BotMessage{ChatId: chat_id}
		if count == 0 {
			botMessage.Text = "You have done nothing:( Maybe you want to get good result"
		} else if count >= 1 && count <= 5 {
			botMessage.Text = "Not Bad, today you have done " + strconv.Itoa(count) + " doings, try better"
		} else {
			botMessage.Text = "Great, today you have done " + strconv.Itoa(count) + " doings, Keep it"
		}
		language := postgresql.GetLanguage(botMessage.ChatId)
		if language == "Russian" {
			botMessage.Text = parser.Translate(botMessage.Text)
		}
		buf, _ := json.Marshal(botMessage)
		_, err := http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Println("Hi")
			log.Fatal(err)
		}
	}
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].count < users[j].count
	})
	for _, use := range users {
		percent := use.count * 100 / len(users)
		buf, _ := json.Marshal(model.BotMessage{ChatId: use.chat_id, Text: "Today you are better than " + strconv.Itoa(min(100, percent)) + "% of users."})
		_, err := http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Println("Hi")
			log.Fatal(err)
		}
	}
}
