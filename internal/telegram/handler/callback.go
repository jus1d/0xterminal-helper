package handler

import (
	"fmt"
	"strconv"
	"strings"
	"terminal/internal/storage"
	"terminal/internal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CallbackContinueGame(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	messageID := u.CallbackQuery.Message.MessageID

	game, exists := h.games[userID]
	if !exists {
		h.editMessage(userID, messageID, "<b>You have no started games</b>\n\nUse /newgame or button to start new one", GetMarkupNewGame())
		return
	}

	h.editMessage(userID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(game.AvailableWords))
}

func (h *Handler) CallbackStartNewGame(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	h.stages[userID] = WaitingWordList
	delete(h.games, userID)
	h.sendTextMessage(userID, "Send me list of words in your $TERMINAL game", nil)
}

func (h *Handler) CallbackWordsList(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	game := h.games[userID]
	h.editMessage(userID, u.CallbackQuery.Message.MessageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(game.AvailableWords))
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
		h.editMessage(userID, messageID, "Use /newgame or button to start new game", GetMarkupNewGame())
		return
	}

	attempt := terminal.Attempt{
		Word:           word,
		GuessedLetters: guessedLetters,
	}
	game.CommitAttempt(attempt)

	if len(game.AvailableWords) == 1 {
		delete(h.games, userID)
		h.editMessage(userID, messageID, fmt.Sprintf("<b>Target word:</b> <code>%s</code>", game.AvailableWords[0]), GetMarkupNewGame())

		// we'll assume that game is kinda spam, if initial words is less than 6
		if len(game.InitialWords) >= 6 {
			storage.SaveGame(storage.ConvertToGame(game, u.CallbackQuery.From.UserName, userID))
		}
		return
	}
	if len(game.AvailableWords) == 0 {
		delete(h.games, userID)
		h.editMessage(userID, messageID, "<b>No matching words left.</b>\n\nTry again, may be you made a mistake?", nil)
		return
	}

	h.editMessage(userID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(h.games[userID].AvailableWords))
}
