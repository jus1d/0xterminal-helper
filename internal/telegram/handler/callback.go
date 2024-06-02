package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"terminal/internal/storage"
	"terminal/internal/terminal"
	"terminal/pkg/log/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CallbackContinueGame(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID

	game, exists := h.games[author.ID]
	if !exists {
		h.editMessage(author.ID, messageID, "<b>You have no started games</b>\n\nUse /newgame or button to start new one", GetMarkupNewGame())
		return
	}

	h.editMessage(author.ID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(game.AvailableWords))
}

func (h *Handler) CallbackStartNewGame(u tgbotapi.Update) {
	author := u.CallbackQuery.From

	h.stages[author.ID] = WaitingWordList
	delete(h.games, author.ID)
	h.sendTextMessage(author.ID, "Send me list of words in your $TERMINAL game", nil)
}

func (h *Handler) CallbackWordsList(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID

	game, exists := h.games[author.ID]
	if !exists {
		h.editMessage(author.ID, messageID, "<b>Use /newgame or button to start new game</b>", GetMarkupNewGame())
		return
	}
	h.editMessage(author.ID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(game.AvailableWords))
}

func (h *Handler) CallbackChooseWord(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID

	parts := strings.Split(u.CallbackData(), ":")
	word := parts[1]
	h.editMessage(author.ID, messageID, fmt.Sprintf("<b>How many guessed letters in word</b> <code>%s</code>?", word), GetMarkupGuessedLetters(word))
}

func (h *Handler) CallbackChooseGuessedLetters(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	log := h.log.With(
		slog.String("op", "handler.CallbackChooseGuessedLetters"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
		slog.String("query", u.CallbackData()),
	)

	messageID := u.CallbackQuery.Message.MessageID

	_, err := h.storage.GetUserByTelegramID(author.ID)
	if errors.Is(err, storage.ErrUserNotFound) {
		_, err = h.storage.SaveUser(author.ID, author.UserName, author.FirstName, author.LastName)
		if err != nil {
			log.Error("could not save user to database", sl.Err(err))
		}
	} else if err != nil {
		log.Error("failed to get user from database", sl.Err(err))
		return
	}

	parts := strings.Split(u.CallbackData(), ":")[1:]

	word := parts[0]
	guessedLetters, _ := strconv.Atoi(parts[1])

	game, exists := h.games[author.ID]
	if !exists {
		h.editMessage(author.ID, messageID, "Use /newgame or button to start new game", GetMarkupNewGame())
		return
	}

	attempt := terminal.Attempt{
		Word:           word,
		GuessedLetters: guessedLetters,
	}
	game.SubmitAttempt(attempt)

	if len(game.AvailableWords) == 1 {
		delete(h.games, author.ID)
		h.editMessage(author.ID, messageID, fmt.Sprintf("<b>Target word:</b> <code>%s</code>", game.AvailableWords[0]), GetMarkupNewGame())

		// we'll assume that game is kinda spam, if initial words is less than 6
		if len(game.InitialWords) >= 6 {
			h.storage.SaveGame(author.ID, game.InitialWords, game.AvailableWords[0])
		}
		return
	}
	if len(game.AvailableWords) == 0 {
		delete(h.games, author.ID)
		h.editMessage(author.ID, messageID, "<b>No matching words left.</b>\n\nTry again, may be you made a mistake?", nil)
		return
	}

	h.editMessage(author.ID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(h.games[author.ID].AvailableWords))
}
