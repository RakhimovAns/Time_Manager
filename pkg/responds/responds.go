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

func SetLanguage(data []string, botMessage *model.BotMessage) {
	language := data[1]
	fmt.Println(language)
	if language != "English" && language != "Russian" {
		botMessage.Text = "invalid language"
		return
	}
	postgresql.SetLanguage(botMessage.ChatId, language)
	botMessage.Text = "Was set successfully"
}
func HelpRespond(config *tgbotapi.MessageConfig) {
	config.Text = model.English["HelpRespond"]
}
func StartRespond(config *tgbotapi.MessageConfig) {
	config.Text = "Hello dear user! This bot sorts your doings, to get more info use command /help"
}
func InfoRespond(config *tgbotapi.MessageConfig) {
	config.Text = "This bot sorts your doing by ABCDE method.\n" + "ABCDE method is the one of the most popular sorting methods of doing.The essence of the technique is to sort tasks by importance  using a special table"
}

func DoneRespond(data []string, botMessage *tgbotapi.MessageConfig) {
	Doings, err := parser.Parse(data, botMessage)
	if err != nil {
		return
	}
	for _, doing := range Doings {
		postgresql.SetStatus(doing)
	}
	botMessage.Text = "Command finished successfully"
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
			botMessage.Text = "Have you done anything from your doing list?"
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
	config.Text = "Ansar Rakhmimov. support: @Rakhimov_Ans"
}

func ErrorRespond(config *tgbotapi.MessageConfig) {
	config.Text = "unrecognized command, use /help to get list of commands"
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
		botMessage.Text = "I will remind about that"
	}
}

func Remind(botURL string) {
	Doings := postgresql.GetDoings()
	for _, doing := range Doings {
		var BotMessage model.BotMessage
		BotMessage.ChatId = doing.ChatId
		if doing.Data.Sub(time.Now().Add(3*time.Hour)) <= 0 {
			postgresql.SetStatus(doing)
			BotMessage.Text = "You need to start '" + doing.Name + "'"
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

func ListRespond(botMessage *tgbotapi.MessageConfig) {
	Doings := postgresql.GetDoingsByID(botMessage.ChatID)
	if len(Doings) == 0 {
		botMessage.Text = "You have no doings"
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
	botMessage.Text = answer
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
	botMessage.Text = "Was deleted successfully"
}

func Want(botURL string) {
	if !(time.Now().Hour() == 23 && time.Now().Minute() == 14) {
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
