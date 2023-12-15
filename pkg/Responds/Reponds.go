package Responds

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

func GetPool() *pgxpool.Pool {
	return pool
}

func ConnectToDB(dsn string) {
	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pool, _ = pgxpool.Connect(connectCtx, dsn)
}

func GetUpdates(BotURL string, offset int) ([]types.Update, error) {
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

func Respond(botURL string, update types.Update) error {
	botMessage := new(types.BotMessage)
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
		"/sort - sorts your doings, use this command like following format:\n" + //implemented
		"	Name Date Time Importance(from 1 to 4, from lower to higher)" +
		"/remind - reminds you about your doing, use this command like sort command\n" + //implemented
		"/author - gets information about authors\n" + //implemented //implemented
		"/delete - deletes doing from remind list,use this command like sort command\n" +
		"/list - gets all doing from remind list" // correct english grammar //implemented
	//"/change" //Add it later
}
func StartRespond(botMessage *types.BotMessage) {
	botMessage.Text = "Hello dear user! This bot sorts your doings, to get more info use command /help"
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
	for _, doing := range Doings {
		_, err := pool.Exec(context.Background(), `
				insert into doings(chat_id, name, importance, time) values ($1,$2,$3,$4)
`, botMessage.ChatId, doing.Name, doing.Importance, doing.Data)
		if err != nil {
			log.Fatal("error with adding to db:", err)
		}
	}
	botMessage.Text = "I will remind about it"
}

func GetDoings() []types.DoWithID {
	rows, err := pool.Query(context.Background(), "SELECT id,doings.chat_id,name, time, importance FROM doings")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doings []types.DoWithID

	for rows.Next() {
		var ID int
		var ChatID int
		var name string
		var timestamp time.Time
		var importance int

		err = rows.Scan(&ID, &ChatID, &name, &timestamp, &importance)
		if err != nil {
			log.Fatal(err)
		}

		doing := types.DoWithID{
			ID:         ID,
			ChatId:     ChatID,
			Name:       name,
			Data:       timestamp,
			Importance: importance,
		}
		doings = append(doings, doing)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return doings
}

func Remind(botURL string) {
	Doings := GetDoings()
	for _, doing := range Doings {

		var BotMessage types.BotMessage
		BotMessage.ChatId = doing.ChatId

		if doing.Data.Sub(time.Now().Add(3*time.Hour)).Minutes() == 5 {
			BotMessage.Text = "You need to start " + doing.Name + " after 5 minutes"
			buf, err := json.Marshal(BotMessage)
			if err != nil {
				log.Fatal(err)
			}
			_, err = http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
			if err != nil {
				log.Fatal(err)
			}
		} else if doing.Data.Sub(time.Now().Add(3*time.Hour)) <= 0 {
			BotMessage.Text = "You need to start " + doing.Name
			Delete(doing)
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
func Delete(doing types.DoWithID) {
	_, err := pool.Exec(context.Background(), `
		DELETE FROM doings where id=$1
`, doing.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func ListRespond(botMessage *types.BotMessage) {
	Doings := GetDoingsByID(botMessage.ChatId)
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

func GetDoingsByID(ID int) []types.DoWithID {
	rows, err := pool.Query(context.Background(), "SELECT id,doings.chat_id,name, time, importance FROM doings where chat_id=$1", ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doings []types.DoWithID

	for rows.Next() {
		var ID int
		var ChatID int
		var name string
		var timestamp time.Time
		var importance int
		err = rows.Scan(&ID, &ChatID, &name, &timestamp, &importance)
		if err != nil {
			log.Fatal(err)
		}

		doing := types.DoWithID{
			ID:         ID,
			ChatId:     ChatID,
			Name:       name,
			Data:       timestamp,
			Importance: importance,
		}
		doings = append(doings, doing)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return doings
}
func DeleteRespond(data []string, botMessage *types.BotMessage) {
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
	for _, doing := range Doings {
		_, err := pool.Exec(context.Background(), `
		DELETE FROM doings where name=$1 and importance=$2 and time=$3
`, doing.Name, doing.Importance, doing.Data)
		if err != nil {
			log.Fatal(err)
		}
	}
	botMessage.Text = "was deleted successfully"
}
