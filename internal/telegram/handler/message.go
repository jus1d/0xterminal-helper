package handler

import (
	"errors"
	"strings"
	"terminal/internal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) TextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID
	stage, exists := h.stages[userID]
	if !exists {
		h.stages[userID] = None
	}
	stage, _ = h.stages[userID]

	switch stage {
	case WaitingWordList:
		words := terminal.RemoveTrashFromWordsList(strings.Split(u.Message.Text, "\n"))

		if len(words) < 6 {
			h.sendTextMessage(userID, "<b>Word list can't be too short</b>\n\nSend me list of at least 6 unique words", nil)
			return
		}

		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			h.sendTextMessage(userID, "Words must be the same length", nil)
			return
		}
		h.games[userID] = game
		h.stages[userID] = None

		h.sendTextMessage(userID, "Choose picked word below", GetMarkupWords(h.games[userID].AvailableWords))
	case None:
		h.sendTextMessage(userID, "Use /newgame or click the button to start new $TERMINAL game", GetMarkupNewGame())
	}
}
