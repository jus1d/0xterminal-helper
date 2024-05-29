package handler

import (
	"log/slog"
	"terminal/internal/storage"
	"terminal/internal/terminal"
	"terminal/pkg/log/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stage uint8

const (
	None = iota
	WaitingWordList
)

type Handler struct {
	log     *slog.Logger
	client  *tgbotapi.BotAPI
	storage storage.Storage
	games   map[int64]*terminal.Game
	stages  map[int64]Stage
}

func New(logger *slog.Logger, client *tgbotapi.BotAPI, st storage.Storage) *Handler {
	return &Handler{
		log:     logger,
		client:  client,
		storage: st,
		games:   make(map[int64]*terminal.Game, 0),
		stages:  make(map[int64]Stage, 0),
	}
}

func (h *Handler) sendTextMessage(chatID int64, content string, markup *tgbotapi.InlineKeyboardMarkup) {
	log := h.log.With(
		slog.String("op", "handler.sendTextMessage"),
	)

	message := tgbotapi.NewMessage(chatID, content)
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup

	_, err := h.client.Send(message)
	if err != nil {
		log.Error("could not send message", sl.Err(err))
	}
}

func (h *Handler) editMessage(chatID int64, messageID int, content string, markup *tgbotapi.InlineKeyboardMarkup) {
	log := h.log.With(
		slog.String("op", "handler.editMessage"),
	)

	message := tgbotapi.NewEditMessageText(chatID, messageID, content)
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup

	_, err := h.client.Send(message)
	if err != nil {
		log.Error("could not send message", sl.Err(err))
	}
}
