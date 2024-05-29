package postgres

import (
	"fmt"
	"terminal/internal/config"
	"terminal/pkg/log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

const (
	StageNone = iota
	StageWaintgWordList
)

type User struct {
	ID         string    `db:"id"`
	TelegramID int64     `db:"telegram_id"`
	Username   string    `db:"username"`
	FirstName  string    `db:"firstname"`
	LastName   string    `db:"lastname"`
	Stage      int       `db:"stage"`
	IsAdmin    bool      `db:"is_admin"`
	CreatedAt  time.Time `db:"created_at"`
}

type Game struct {
	ID         string    `db:"id"`
	TelegramID int64     `db:"telegram_id"`
	Words      []string  `db:"words"`
	Target     string    `db:"target"`
	WordsHash  string    `db:"words_hash"`
	Hash       string    `db:"hash"`
	CreatedAt  time.Time `db:"created_at"`
}

func New(conf config.Postgres) *Storage {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Name, conf.Password, conf.ModeSSL))
	if err != nil {
		log.Fatal("could not start postgres database", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("no response from postgres database", err)
	}

	return &Storage{
		db: db,
	}
}
