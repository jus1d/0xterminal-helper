package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terminal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) TextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID
	stage, exists := h.stages[userID]
	if exists && stage == WaitingWordList {
		words := strings.Split(u.Message.Text, "\n")
		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			h.sendTextMessage(userID, "Words must be the same length")
			return
		}
		h.games[userID] = game
		h.stages[userID] = WaitingAttempt

		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(h.games[userID].Words))
		for i, word := range h.games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		h.sendTextMessage(userID, content)
	} else if exists && stage == WaitingAttempt {
		parts := strings.Split(u.Message.Text, " ")
		if len(parts) < 2 {
			h.sendTextMessage(userID, "U invalid. Use: word guessed-letters-amount")
			return
		}
		word := parts[0]
		guessedLetters, err := strconv.Atoi(parts[1])
		if err != nil {
			h.sendTextMessage(userID, "Guessed letters amount should be integer. Use: word guessed-letters-amount")
			return
		}
		attempt := terminal.Attempt{
			Word:           word,
			GuessedLetters: guessedLetters,
		}
		h.games[userID].CommitAttempt(attempt)
		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(h.games[userID].Words))
		for i, word := range h.games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		h.sendTextMessage(userID, content)
		if len(h.games[userID].Words) <= 1 {
			h.stages[userID] = None
		}
	}
}
