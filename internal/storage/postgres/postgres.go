package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"terminal/internal/config"
	"terminal/internal/storage"
	"terminal/internal/terminal"
	"terminal/pkg/log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
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

func (s *Storage) CreateUser(telegramID int64, username string, firstname string, lastname string) (*storage.User, error) {
	query := "INSERT INTO users (telegram_id, username, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING *"
	row := s.db.QueryRow(query, telegramID, username, firstname, lastname)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user storage.User
	err := row.Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.IsAdmin, &user.CreatedAt)

	return &user, err
}

func (s *Storage) GetUserByTelegramID(telegramID int64) (*storage.User, error) {
	var user storage.User
	err := s.db.QueryRow("SELECT * FROM users WHERE telegram_id = $1", telegramID).Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) SaveGame(telegramID int64, words []string, target string) (*storage.Game, error) {
	query := "INSERT INTO games (telegram_id, words, target, words_hash) VALUES ($1, $2, $3, $4) RETURNING *"
	wordsHash := terminal.ComputeWordsHash(words)

	row := s.db.QueryRow(query, telegramID, pq.Array(words), target, wordsHash)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var game storage.Game
	err := row.Scan(&game.ID, &game.TelegramID, &game.Words, &game.Target, &game.WordsHash, &game.CreatedAt)

	return &game, err
}

func (s *Storage) TryFindAnswer(words []string) string {
	wordsHash := terminal.ComputeWordsHash(words)

	query := "SELECT target FROM games WHERE words_hash = $1"

	var target string
	err := s.db.QueryRow(query, wordsHash).Scan(&target)
	if !errors.Is(err, sql.ErrNoRows) {
		log.Error("could not get game", err)
		return ""
	}

	return target
}
