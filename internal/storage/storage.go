package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"terminal/pkg/log"
)

const Path = "./storage/data.json"

type Data struct {
	Games []Game `json:"games"`
}

type Game struct {
	Words  []string `json:"words"`
	Target string   `json:"target"`
}

func SaveGame(game Game) {
	jsonData, err := os.ReadFile(Path)
	if err != nil {
		log.Error("could not read JSON file", err, log.WithString("path", Path))
		return
	}

	var data Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("could not unmarshall JSON data", err)
		return
	}

	data.Games = append(data.Games, game)
	SaveData(&data)
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
