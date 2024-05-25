package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	userID := u.Message.From.ID
	h.sendTextMessage(userID, "Use /newgame to register a game", nil)
	h.stages[userID] = None
}

func (h *Handler) CommandGame(u tgbotapi.Update) {
	userID := u.Message.From.ID
	game, exists := h.games[userID]
	if exists {
		content := "<b>You already have started game. Do you want to continue it?</b>\n\n<b>Words:</b>\n"
		for i, word := range game.Words {
			content += fmt.Sprintf("<code>%s</code>\n", i+1, word)
		}
		h.sendTextMessage(userID, content, GetMarkupGameMenu())
	} else {
		h.sendTextMessage(userID, "Send me list of words in your $TERMINAL game", nil)
		h.stages[userID] = WaitingWordList
	}
}
