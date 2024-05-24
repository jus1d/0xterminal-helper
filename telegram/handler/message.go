package handler

import (
	"errors"
	"strings"
	"terminal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) TextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID
	stage, exists := h.stages[userID]
	if !exists {
		h.sendTextMessage(userID, "Use /start to start the bot", nil)
		return
	}

	switch stage {
	case WaitingWordList:
		words := strings.Split(u.Message.Text, "\n")
		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			h.sendTextMessage(userID, "Words must be the same length", nil)
			return
		}
		h.games[userID] = game
		h.stages[userID] = None

		h.sendTextMessage(userID, "Choose picked word below", GetMarkupWords(h.games[userID].Words))
	}
}
