package handler

import (
	"fmt"
	"strconv"
	"strings"
	"terminal/internal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CallbackContinueGame(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	messageID := u.CallbackQuery.Message.MessageID

	game, exists := h.games[userID]
	if !exists {
		h.editMessage(userID, messageID, "<b>You have no started games</b>\nUse /newgame to start new one", nil)
		return
	}

	h.editMessage(userID, messageID, "Choose picked word below", GetMarkupWords(game.Words))
}

func (h *Handler) CallbackStartNewGame(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	h.stages[userID] = WaitingWordList
	h.editMessage(userID, u.CallbackQuery.Message.MessageID, "Send me list of words in your $TERMINAL game", nil)
}

func (h *Handler) CallbackChooseWord(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	parts := strings.Split(u.CallbackData(), ":")
	word := parts[1]
	h.editMessage(userID, u.CallbackQuery.Message.MessageID, fmt.Sprintf("<b>How many guessed letters in word</b> <code>%s</code>?", word), GetMarkupGuessedLetters(word))
}

func (h *Handler) CallbackChooseGuessedLetters(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	messageID := u.CallbackQuery.Message.MessageID
	parts := strings.Split(u.CallbackData(), ":")[1:]

	word := parts[0]
	guessedLetters, _ := strconv.Atoi(parts[1])

	game, exists := h.games[userID]
	if !exists {
		h.editMessage(userID, messageID, "Use /newgame to start ne game", nil)
		return
	}

	attempt := terminal.Attempt{
		Word:           word,
		GuessedLetters: guessedLetters,
	}
	game.CommitAttempt(attempt)

	if len(game.Words) == 1 {
		delete(h.games, userID)
		h.editMessage(userID, messageID, fmt.Sprintf("<b>Target word:</b> <code>%s</code>", game.Words[0]), nil)
		return
	}
	if len(game.Words) == 0 {
		delete(h.games, userID)
		h.editMessage(userID, messageID, "<b>No possible words left.</b>\nTry again, may be you made a mistake?", nil)
		return
	}

	h.editMessage(userID, messageID, "Choose picked word below", GetMarkupWords(h.games[userID].Words))
}
