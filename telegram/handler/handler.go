package handler

import (
	"log"
	"terminal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stage uint8

const (
	None = iota
	WaitingWordList
)

type Handler struct {
	client *tgbotapi.BotAPI
	games  map[int64]*terminal.Game
	stages map[int64]Stage
}

func New(client *tgbotapi.BotAPI) *Handler {
	return &Handler{
		client: client,
		games:  make(map[int64]*terminal.Game, 0),
		stages: make(map[int64]Stage, 0),
	}
}

func (h *Handler) sendTextMessage(chatID int64, content string, markup *tgbotapi.InlineKeyboardMarkup) {
	message := tgbotapi.NewMessage(chatID, content)
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup

	_, err := h.client.Send(message)
	if err != nil {
		log.Printf("ERROR: could not send message to ID: %d, error: %s\n", chatID, err.Error())
		return
	}
}

func (h *Handler) editMessage(chatID int64, messageID int, content string, markup *tgbotapi.InlineKeyboardMarkup) {
	message := tgbotapi.NewEditMessageText(chatID, messageID, content)
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup

	_, err := h.client.Send(message)
	if err != nil {
		log.Printf("ERROR: could not send message to ID: %d, error: %s\n", chatID, err.Error())
		return
	}
}
