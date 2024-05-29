package main

import (
	"encoding/json"
	"os"
	"terminal/internal/config"
	"terminal/internal/storage/postgres"
	"terminal/pkg/log"
)

type Data struct {
	TotalGames int    `json:"total_games"`
	Games      []Game `json:"games"`
}

type Game struct {
	Words     []string `json:"words"`
	Target    string   `json:"target"`
	WordsHash string   `json:"words_hash"`
	PlayedBy  User     `json:"played_by"`
}

type User struct {
	Username   string `json:"username"`
	TelegramID int64  `json:"telegram_id"`
}

func main() {
	conf := config.MustLoad()
	storage := postgres.New(conf.Postgres)

	data := loadData()

	for _, game := range data.Games {
		storage.SaveGame(game.PlayedBy.TelegramID, game.Words, game.Target)
	}
}

func loadData() *Data {
	path := "./storage/data.json"
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Error("could not read JSON file", err, log.WithString("path", path))
		return nil
	}

	var data Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("could not unmarshall JSON data", err)
		return nil
	}

	return &data
}
