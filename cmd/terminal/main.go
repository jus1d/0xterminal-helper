package main

import (
	"os"
	"terminal/internal/config"
	"terminal/internal/ocr"
	"terminal/internal/storage/postgres"
	"terminal/internal/telegram"
	"terminal/pkg/log"
	"terminal/pkg/log/sl"

	"github.com/robfig/cron/v3"
)

func main() {
	conf := config.MustLoad()

	logger := log.Init(conf.Env)

	if conf.Env == config.EnvProduction {
		c := cron.New()

		c.AddFunc("0 0 * * *", func() { // update logger every day in 00:00
			*logger = *log.Init(conf.Env)
		})

		c.Start()
	}

	storage, err := postgres.New(conf.Postgres)
	if err != nil {
		logger.Error("failed to connect to postgres database", sl.Err(err))
		os.Exit(1)
	}

	bot := telegram.New(logger, conf.Telegram, storage, ocr.New(conf.OCR.Tokens))
	bot.Run()
}
