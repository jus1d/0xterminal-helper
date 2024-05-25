package main

import (
	"os"
	"terminal/internal/storage"
	"terminal/internal/telegram"
	"terminal/pkg/log"

	"github.com/joho/godotenv"
)

func main() {
	initStorage()

	if err := godotenv.Load(); err != nil {
		log.Fatal("could not load `.env` file", err)
	}

	token := os.Getenv("TOKEN")
	bot := telegram.New(token)
	bot.Run()
}

func initStorage() {
	if _, err := os.Stat(storage.Path); os.IsNotExist(err) {
		err = os.Mkdir("./storage", 0755)
		if err != nil {
			log.Fatal("could not create storage folder", err)
		}
	}

	if _, err := os.Stat(storage.Path); os.IsNotExist(err) {
		_, err = os.Create(storage.Path)
		if err != nil {
			log.Fatal("could not create storage file", err)
		}
		data := &storage.Data{
			Games: make([]storage.Game, 0),
		}
		storage.SaveData(data)
	}
}
