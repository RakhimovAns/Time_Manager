package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/RakhimovAns/Time_Manager/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var pool *pgxpool.Pool

func main() {
	BotToken := "6791006120:AAFzk656CBCPbNlWVolFZl1cUwp4ej7A-Tc"
	//https://api.telegram.org/bot<token>/METHOD_NAME
	BotAPI := "https://api.telegram.org/bot"
	BotURL := BotAPI + BotToken
	offset := 0
	dsn := "postgresql://postgres:postgres@localhost:5432/manager"
	ConnectToDB(dsn)
	defer pool.Close()
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
func ConnectToDB(dsn string) {
	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pool, _ = pgxpool.Connect(connectCtx, dsn)
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
	data := strings.Split(strings.Replace(update.Message.Text, "\r\n", "\n", -1), "\n")
	if update.Message.Text == "/help" {
		HelpRespond(botMessage)
	} else if update.Message.Text == "/author" {
		AuthorRespond(botMessage)
	} else if update.Message.Text == "/info" {
		InfoRespond(botMessage)
	} else if data[0] == "/sort" && len(data) > 1 {
		SortRespond(data[1:], botMessage)
	} else if data[0] == "/remind" && len(data) > 1 {
		RemindRespond(data[1:], botMessage)
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
		"/sort - sorts your doings, use this command like following format:\n" +
		"	Name Date Time Importance(from 1 to 4, from lower to higher)" +
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

func SortRespond(data []string, botMessage *types.BotMessage) {
	var Doings []types.Doing
	for _, doing := range data {
		SplitedData := strings.Split(doing, " ")
		if len(SplitedData) != 4 {
			botMessage.Text = "invalid type of doings"
			return
		}
		var Do types.Doing
		Do.Name = SplitedData[0]
		DateTimeStr := SplitedData[1] + " " + SplitedData[2]
		layout := "2.01.2006 15:04"
		dateTime, err := time.Parse(layout, DateTimeStr)
		if err != nil {
			botMessage.Text = "invalid type of doing"
			return
		}
		Do.Data = dateTime
		Do.Importance, _ = strconv.Atoi(SplitedData[3])
		Doings = append(Doings, Do)
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
func RemindRespond(data []string, botMessage *types.BotMessage) {
	var Doings []types.Doing
	for _, doing := range data {
		SplitedData := strings.Split(doing, " ")
		if len(SplitedData) != 4 {
			botMessage.Text = "invalid type of doings"
			return
		}
		var Do types.Doing
		Do.Name = SplitedData[0]
		DateTimeStr := SplitedData[1] + " " + SplitedData[2]
		layout := "2.01.2006 15:04"
		dateTime, err := time.Parse(layout, DateTimeStr)
		if err != nil {
			botMessage.Text = "invalid type of doing"
			return
		}
		Do.Data = dateTime
		Do.Importance, _ = strconv.Atoi(SplitedData[3])
		Doings = append(Doings, Do)
	}
	sort.SliceStable(Doings, func(i, j int) bool {
		if Doings[i].Data != Doings[j].Data {
			return Doings[i].Data.Before(Doings[j].Data)
		}
		return Doings[i].Importance > Doings[j].Importance
	})
	ctx := context.Background()
	for _, doing := range Doings {
		_, err := pool.Exec(ctx, `
				insert into doings(chat_id, name, importance, time) values ($1,$2,$3,$4)
`, botMessage.ChatId, doing.Name, doing.Importance, doing.Data)
		if err != nil {
			log.Fatal("error with adding to db:", err)
		}
	}
	botMessage.Text = "I will remind about it"
}

func Remind() {

}
