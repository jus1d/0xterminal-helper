package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
