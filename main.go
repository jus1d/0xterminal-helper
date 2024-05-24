package main

import (
	"fmt"
	"os"
	"terminal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("ERROR: could not load .env file: %v", err)
		return
	}

	token := os.Getenv("TOKEN")
	bot := telegram.New(token)
	bot.Run()
}
