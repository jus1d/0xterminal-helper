package main

import (
	"os"
	"terminal/internal/telegram"
	"terminal/pkg/log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("could not load `.env` file", err)
	}

	token := os.Getenv("TOKEN")
	bot := telegram.New(token)
	bot.Run()
}
