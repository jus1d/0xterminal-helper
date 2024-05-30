package main

import (
	"os"
	"terminal/internal/config"
	"terminal/internal/ocr"
	"terminal/internal/storage/postgres"
	"terminal/internal/telegram"
	"terminal/pkg/log"
	"terminal/pkg/log/sl"
)

func main() {
	conf := config.MustLoad()

	logger := log.Init(conf.Env)

	storage, err := postgres.New(conf.Postgres)
	if err != nil {
		logger.Error("failed to start postgres database", sl.Err(err))
		os.Exit(1)
	}

	bot := telegram.New(logger, conf.Telegram, storage, ocr.New(conf.OCR.Token))
	bot.Run()
}
