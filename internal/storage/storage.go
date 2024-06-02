package storage

import (
	"terminal/internal/terminal/dataset"
	"time"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

// TODO(#9): add custom errors to storage functions

var (
	ErrUserNotFound = errors.New("0xterminal.storage: user not found")
)

type Storage interface {
	CreateUser(telegramID int64, username string, firstname string, lastname string) (*User, error)
	GetUserByTelegramID(telegramID int64) (*User, error)
	SaveGame(telegramID int64, words []string, target string) (*Game, error)
	TryFindAnswer(words []string) (string, error)
	GetDataset() (*dataset.Dataset, error)
	GetDailyReport() (*DailyReport, error)
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
	IsAdmin    bool      `db:"is_admin"`
	CreatedAt  time.Time `db:"created_at"`
}

type Game struct {
	ID         string    `db:"id" json:""`
	TelegramID int64     `db:"telegram_id"`
	Words      []string  `db:"words"`
	Target     string    `db:"target"`
	WordsHash  string    `db:"words_hash"`
	CreatedAt  time.Time `db:"created_at"`
}

type DailyReport struct {
	Stats       []DailyUserStat
	JoinedUsers []string
}

type DailyUserStat struct {
	Username    string
	GamesPlayed int
}
