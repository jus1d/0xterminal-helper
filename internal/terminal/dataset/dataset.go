package dataset

import (
	"encoding/json"
	"os"
	"time"
)

type Dataset struct {
	TotalGames int    `json:"total_games"`
	Games      []Game `json:"games"`
}

type Game struct {
	Words          []string  `json:"words"`
	Target         string    `json:"target"`
	AttemptsAmount int       `json:"attempts_amount"`
	User           User      `json:"user"`
	WordsHash      string    `json:"words_hash"`
	CreatedAt      time.Time `json:"created_at"`
}

type User struct {
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
}

func ExportDatasetToJSON(data *Dataset) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	path := time.Now().Format("dataset-02-01-2006.json")
	err = os.WriteFile(path, jsonData, 0600)
	if err != nil {
		return "", err
	}

	return path, nil
}
