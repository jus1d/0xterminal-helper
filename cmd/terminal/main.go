package main

import (
	"terminal/internal/config"
	"terminal/internal/storage/postgres"
	"terminal/internal/telegram"
)

func main() {
	conf := config.MustLoad()

	storage := postgres.New(conf.Postgres)

	bot := telegram.New(conf.Telegram, storage)
	bot.Run()
}
