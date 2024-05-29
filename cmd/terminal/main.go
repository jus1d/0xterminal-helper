package main

import (
	"os"
	"terminal/internal/config"
	"terminal/internal/storage"
	"terminal/internal/telegram"
	"terminal/pkg/log"
)

func main() {
	initStorage()
	conf := config.MustLoad()

	bot := telegram.New(conf.Telegram)
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
