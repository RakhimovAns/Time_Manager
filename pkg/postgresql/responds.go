package postgresql

import (
	"context"
	"github.com/RakhimovAns/Time_Manager/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

var pool *pgxpool.Pool

func GetPool() *pgxpool.Pool {
	return pool
}

func ConnectToDB(dsn string) {
	pool, _ = pgxpool.Connect(context.Background(), dsn)

}

func GetDoingsWithStatus() []model.Doing {
	rows, err := pool.Query(context.Background(), "SELECT id,doings.chat_id,name, time, importance FROM doings where status=false")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doings []model.Doing

	for rows.Next() {
		var ID int
		var ChatID int64
		var name string
		var timestamp time.Time
		var importance int

		err = rows.Scan(&ID, &ChatID, &name, &timestamp, &importance)
		if err != nil {
			log.Fatal(err)
		}

		doing := model.Doing{
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

func GetDoingsByID(ID int64) []model.Doing {
	rows, err := pool.Query(context.Background(), "SELECT id,doings.chat_id,name, time, importance FROM doings where chat_id=$1 and status=false", ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doings []model.Doing

	for rows.Next() {
		var ID int
		var ChatID int64
		var name string
		var timestamp time.Time
		var importance int
		err = rows.Scan(&ID, &ChatID, &name, &timestamp, &importance)
		if err != nil {
			log.Fatal(err)
		}

		doing := model.Doing{
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
func SetLanguage(ID int64, language string) {
	_, err := pool.Exec(context.Background(), `
	 insert into languages(chat_id, language) values ($1,$2)
`, ID, language)
	if err != nil {
		log.Println("Error to set right language")
		log.Fatal(err)
	}
}
func SetStatus(doing model.Doing) {
	_, err := pool.Exec(context.Background(), `
		update doings set status=true where name=$1 and time=$2 and importance=$3
`, doing.Name, doing.Data, doing.Importance)
	if err != nil {
		log.Fatal(err)
	}
}
func Delete(doing model.Doing) error {
	_, err := pool.Exec(context.Background(), `
		DELETE FROM doings where name=$1 and time=$2 and importance=$3
`, doing.Name, doing.Data, doing.Importance)
	if err != nil {
		return err
	}
	return nil
}

func GetAllDoings() []model.Doing {
	rows, err := pool.Query(context.Background(), "SELECT id,chat_id,name, time, importance,status FROM doings")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doings []model.Doing

	for rows.Next() {
		var ID int
		var ChatID int64
		var name string
		var timestamp time.Time
		var importance int
		var status bool
		err = rows.Scan(&ID, &ChatID, &name, &timestamp, &importance, &status)
		if err != nil {
			log.Fatal(err)
		}

		doing := model.Doing{
			ID:         ID,
			ChatId:     ChatID,
			Name:       name,
			Data:       timestamp,
			Importance: importance,
			Status:     status,
		}
		doings = append(doings, doing)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return doings
}

func GetDoings() []model.Doing {
	rows, err := pool.Query(context.Background(), "SELECT id,doings.chat_id,name, time, importance,status FROM doings where status=false")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var doings []model.Doing

	for rows.Next() {
		var ID int
		var ChatID int64
		var name string
		var timestamp time.Time
		var importance int
		var status bool
		err = rows.Scan(&ID, &ChatID, &name, &timestamp, &importance, &status)
		if err != nil {
			log.Fatal(err)
		}

		doing := model.Doing{
			ID:         ID,
			ChatId:     ChatID,
			Name:       name,
			Data:       timestamp,
			Importance: importance,
			Status:     status,
		}
		doings = append(doings, doing)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return doings
}
func AddDoings(doing model.Doing, botMessage *tgbotapi.MessageConfig) error {
	_, err := pool.Exec(context.Background(), `
			insert into doings(chat_id, name, importance, time) values ($1,$2,$3,$4)
			`, botMessage.ChatID, doing.Name, doing.Importance, doing.Data)
	if err != nil {
		return err
	}
	return nil
}
