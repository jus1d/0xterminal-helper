package handler

import (
	"log/slog"
	"terminal/internal/ocr"
	"terminal/internal/storage"
	"terminal/internal/terminal"
	"terminal/pkg/log/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	GreetingSticker = tgbotapi.FileID("CAACAgIAAxkBAAIEZGZgUk4PpvRpPIEzmIF5SnLlRPsCAAJ1WgAC0YohSqBRt93rOG5hNQQ")
	WaitingSticker  = tgbotapi.FileID("CAACAgIAAxkBAAIEYGZgG0yU3WUeIN7d_brzaqUEchPtAAIaSQACsCNJSmO4cga8SZwHNQQ")
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
	ocr     *ocr.Client
	games   map[int64]*terminal.Game
	stages  map[int64]Stage
}

func New(logger *slog.Logger, client *tgbotapi.BotAPI, st storage.Storage, o *ocr.Client) *Handler {
	return &Handler{
		log:     logger,
		client:  client,
		storage: st,
		ocr:     o,
		games:   make(map[int64]*terminal.Game, 0),
		stages:  make(map[int64]Stage, 0),
	}
}

func (h *Handler) sendTextMessage(chatID int64, content string, markup *tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	log := h.log.With(
		slog.String("op", "handler.sendTextMessage"),
	)

	chattable := tgbotapi.NewMessage(chatID, content)
	chattable.ParseMode = tgbotapi.ModeHTML
	chattable.ReplyMarkup = markup

	message, err := h.client.Send(chattable)
	if err != nil {
		log.Error("could not send message", sl.Err(err))
	}
	return message, err
}

func (h *Handler) editMessage(chatID int64, messageID int, content string, markup *tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	log := h.log.With(
		slog.String("op", "handler.editMessage"),
	)

	chattable := tgbotapi.NewEditMessageText(chatID, messageID, content)
	chattable.ParseMode = tgbotapi.ModeHTML
	chattable.ReplyMarkup = markup

	message, err := h.client.Send(chattable)
	if err != nil {
		log.Error("could not send message", sl.Err(err))
	}

	return message, err
}

func (h *Handler) sendSticker(chatID int64, sticker tgbotapi.RequestFileData) (tgbotapi.Message, error) {
	log := h.log.With(
		slog.String("op", "handler.sendSticker"),
	)

	chattable := tgbotapi.NewSticker(chatID, sticker)
	message, err := h.client.Send(chattable)
	if err != nil {
		log.Error("could not send sticker message", sl.Err(err))
	}

	return message, err
}

func (h *Handler) deleteMessage(chatID int64, messageID int) {
	log := h.log.With(
		slog.String("op", "handler.deleteMessage"),
	)

	config := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}
	_, err := h.client.Request(config)
	if err != nil {
		log.Error("could not delete message", sl.Err(err))
	}
}
