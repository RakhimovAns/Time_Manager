package postgresql

import (
	"context"
	"github.com/RakhimovAns/Time_Manager/model"
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
func AddDoings(doing model.Doing, botMessage *model.BotMessage) error {
	_, err := pool.Exec(context.Background(), `
			insert into doings(chat_id, name, importance, time) values ($1,$2,$3,$4)
			`, botMessage.ChatId, doing.Name, doing.Importance, doing.Data)
	if err != nil {
		return err
	}
	return nil
}
