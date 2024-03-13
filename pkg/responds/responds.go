package respond

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RakhimovAns/Time_Manager/model"
	"github.com/RakhimovAns/Time_Manager/pkg/parser"
	"github.com/RakhimovAns/Time_Manager/pkg/postgresql"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetUpdates(BotURL string, offset int) ([]model.Update, error) {
	resp, err := http.Get(BotURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponses model.RestResponse

	err = json.Unmarshal(body, &restResponses)

	if err != nil {
		return nil, err
	}

	return restResponses.Result, nil
}

func Respond(botURL string, update model.Update) error {
	botMessage := new(model.BotMessage)

	botMessage.ChatId = update.Message.Chat.ChatId

	data := strings.Split(strings.Replace(update.Message.Text, "\r\n", "\n", -1), "\n")

	if update.Message.Text == "/start" {
		StartRespond(botMessage)
	} else if update.Message.Text == "/help" {
		HelpRespond(botMessage)
	} else if update.Message.Text == "/author" {
		AuthorRespond(botMessage)
	} else if update.Message.Text == "/info" {
		InfoRespond(botMessage)
	} else if update.Message.Text == "/list" {
		ListRespond(botMessage)
	} else if data[0] == "/sort" && len(data) > 1 {
		SortRespond(data[1:], botMessage)
	} else if data[0] == "/remind" && len(data) > 1 {
		RemindRespond(data[1:], botMessage)
	} else if data[0] == "/delete" && len(data) > 1 {
		DeleteRespond(data[1:], botMessage)
	} else if data[0] == "/done" && len(data) > 1 {
		DoneRespond(data[1:], botMessage)
	} else if data[0] == "/lang" && len(data) > 1 {
		SetLanguage(data[0:], botMessage)
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
func HelpRespond(botMessage *model.BotMessage) {

	botMessage.Text = "Hello, this bot can sort your doings and remind about them\n" +
		"You can use this following commands\n" +
		"/info - gets information about sorting methods\n" +
		"/sort - sorts your doings, use this command in following format:\n" +
		"	Name Date Time Importance(from 1 to 4, from lower to higher)\n" +
		"	Example: Task 8.02.2024 13:50 1\n" +
		"/remind - reminds you about your doing, use this command like  a sort command\n" +
		"/author - gets information about authors\n" +
		"/delete - deletes doings from remind list,use this command like a sort command\n" +
		"/list - gets all doings from remind list\n" +
		"/done - you can use this command when you finished some doings, use this command like a sort command"
}
func StartRespond(botMessage *model.BotMessage) {
	botMessage.Text = "Hello dear user! This bot sorts your doings, to get more info use command /help"
}
func InfoRespond(botMessage *model.BotMessage) {
	botMessage.Text = "This bot sorts your doing by ABCDE method.\n" + "ABCDE method is the one of the most popular sorting methods of doing.The essence of the technique is to sort tasks by importance  using a special table"
}
func DoneRespond(data []string, botMessage *model.BotMessage) {
	Doings, err := parser.Pars(data, botMessage)
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
func AuthorRespond(botMessage *model.BotMessage) {
	botMessage.Text = "Ansar Rakhmimov. support: @Rakhimov_Ans"
}

func ErrorRespond(botMessage *model.BotMessage) {
	botMessage.Text = "unrecognized command, use /help to get list of commands"
}

func SortRespond(data []string, botMessage *model.BotMessage) {
	Doings, err := parser.Pars(data, botMessage)
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
func RemindRespond(data []string, botMessage *model.BotMessage) {
	Doings, err := parser.Pars(data, botMessage)
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

func ListRespond(botMessage *model.BotMessage) {
	Doings := postgresql.GetDoingsByID(botMessage.ChatId)
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

func DeleteRespond(data []string, botMessage *model.BotMessage) {
	Doings, err := parser.Pars(data, botMessage)
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
