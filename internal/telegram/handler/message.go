package handler

import (
	"errors"
	"strings"
	"terminal/internal/terminal"
	"terminal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) TextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID

	_, err := h.storage.GetUserByTelegramID(userID)
	if err != nil {
		log.Error("could not get user from database", err, log.WithInt64("telegram_id", userID))
		h.sendTextMessage(userID, "<b>It seems that you are new here</b>\n\nUse /start to start the bot", nil)
		return
	}

	stage, exists := h.stages[userID]
	if !exists {
		h.stages[userID] = None
	}
	stage, _ = h.stages[userID]

	switch stage {
	case WaitingWordList:
		words := terminal.RemoveTrashFromWordList(strings.Split(u.Message.Text, "\n"))

		if len(words) < 6 {
			h.sendTextMessage(userID, "<b>According to the $TERMINAL rules, the word list must consist of at least 6 words</b>\n\nSend me list of words in your $TERMINAL game", nil)
			return
		}

		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			h.sendTextMessage(userID, "<b>According to the $TERMINAL rules, the word list should only consist of words of the same length</b>\n\nSend me list of words in your $TERMINAL game", nil)
			return
		}
		h.games[userID] = game
		h.stages[userID] = None

		// TODO:
		// answer := storage.TryFindAnswer(words)
		// if answer != "" {
		// 	h.sendTextMessage(userID, "<b>Found game with similar words list</b>\n\nProbably, the target is <code>"+answer+"</code>", nil)
		// }

		h.sendTextMessage(userID, "<b>Pick one of the words in the list</b>", GetMarkupWords(h.games[userID].AvailableWords))
	case None:
		h.sendTextMessage(userID, "Use /newgame or click the button to start new $TERMINAL game", GetMarkupNewGame())
	}
}
