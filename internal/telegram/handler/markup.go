package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMarkupGameMenu() *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Continue", "game-continue"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Start new game", "start-new-game"),
		),
	)
	return &markup
}

func GetMarkupWords(words []string) *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, word := range words {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(word, fmt.Sprintf("choose-word:%s", word)),
		))
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup
}

func GetMarkupGuessedLetters(word string) *tgbotapi.InlineKeyboardMarkup {
	rowCapacity := len(word) / 2
	if rowCapacity > 5 {
		rowCapacity = len(word) / 3
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for i := 0; i <= len(word); i += 5 {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		for j := 0; len(row) < 5 && i+j <= len(word); j++ {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+j), fmt.Sprintf("choose-guessed-letters:%s:%d", word, i+j)))
		}
		rows = append(rows, row)
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup
}

func GetMarkupNewGame() *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Start new game", "start-new-game"),
		),
	)
	return &markup
}
