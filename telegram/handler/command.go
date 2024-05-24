package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	userID := u.Message.From.ID
	h.sendTextMessage(userID, "Use /game to register a game")
	h.stages[userID] = None
}

func (h *Handler) CommandGame(u tgbotapi.Update) {
	userID := u.Message.From.ID
	game, exists := h.games[userID]
	if exists {
		content := "<b>You already have started game. Do you want to continue it?</b>\n\nWords:\n"
		for i, word := range game.Words {
			content += fmt.Sprintf("_%d: %s\n", i+1, word)
		}
		h.sendTextMessageWithButtons(userID, content, GetMarkupGameMenu())
	} else {
		h.sendTextMessage(userID, "Send me list of words in your $TERMINAL game")
		h.stages[userID] = WaitingWordList
	}
}
