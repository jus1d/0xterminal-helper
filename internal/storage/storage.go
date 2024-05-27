package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"terminal/internal/terminal"
	"terminal/pkg/log"
)

const Path = "./storage/data.json"

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

func SaveGame(game *Game) {
	data := LoadData()

	data.Games = append(data.Games, *game)
	data.TotalGames = len(data.Games)

	SaveData(data)
}

func SaveData(data *Data) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error("could not marshall data to JSON", err)
		return
	}

	err = os.WriteFile(Path, jsonData, 0644)
	if err != nil {
		log.Error("could not write JSON to file", err, log.WithString("path", Path))
		return
	}

	log.Info(fmt.Sprintf("games data saved to `%s`", Path))
}

func LoadData() *Data {
	jsonData, err := os.ReadFile(Path)
	if err != nil {
		log.Error("could not read JSON file", err, log.WithString("path", Path))
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

func TryFindAnswer(words []string) string {
	data := LoadData()

	hash := terminal.ComputeWordsHash(words)

	for _, game := range data.Games {
		if game.WordsHash == hash {
			return game.Target
		}
	}
	return ""
}

func ConvertToGame(game *terminal.Game, username string, telegramID int64) *Game {
	return &Game{
		Words:  game.InitialWords,
		Target: game.AvailableWords[0],
		PlayedBy: User{
			Username:   username,
			TelegramID: telegramID,
		},
		WordsHash: terminal.ComputeWordsHash(game.InitialWords),
	}
}
