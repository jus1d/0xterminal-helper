package handler

import (
	"errors"
	"fmt"
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
		words := strings.Split(u.Message.Text, "\n")
		if len(words) == 1 {
			h.sendTextMessage(userID, fmt.Sprintf("<b>Target word:</b> <code>%s</code>", words[0]), nil)
			return
		}

		game, err := terminal.New(terminal.RemoveTrashFromWordsList(words))
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
