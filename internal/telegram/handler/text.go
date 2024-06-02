package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"terminal/internal/terminal"
	"terminal/pkg/log/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) TextMessage(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.TextMessage"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	stage, exists := h.stages[author.ID]
	if !exists {
		h.stages[author.ID] = None
	}
	stage, _ = h.stages[author.ID]

	switch stage {
	case WaitingWordList:
		words := terminal.RemoveTrashFromWordList(strings.Split(u.Message.Text, "\n"))

		if len(words) < 6 {
			h.sendTextMessage(author.ID, "<b>According to the $TERMINAL rules, the word list must consist of at least 6 words</b>\n\nSend me list of words in your $TERMINAL game", nil)
			return
		}

		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			h.sendTextMessage(author.ID, "<b>According to the $TERMINAL rules, the word list should only consist of words of the same length</b>\n\nSend me list of words in your $TERMINAL game", nil)
			return
		}
		h.games[author.ID] = game
		h.stages[author.ID] = None

		answer, err := h.storage.TryFindAnswer(words)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("game with sent word list was not found")
			} else {
				log.Error("could not get answer from database", sl.Err(err))
			}
		}
		if answer != "" {
			h.sendTextMessage(author.ID, "<b>Found game with similar words list</b>\n\nProbably, the target is <code>"+answer+"</code>", nil)
		}

		h.sendTextMessage(author.ID, fmt.Sprintf("<b>Pick one of %d words in the list</b>", len(words)), GetMarkupWords(h.games[author.ID].AvailableWords))
	case None:
		h.sendTextMessage(author.ID, "Use /newgame or click the button to start new $TERMINAL game", GetMarkupNewGame())
	}
}
