package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"terminal/internal/storage"
	"terminal/pkg/log/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.CommandStart"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	h.sendSticker(author.ID, GreetingSticker)

	_, err := h.storage.SaveUser(author.ID, author.UserName, author.FirstName, author.LastName)
	if err != nil && !errors.Is(err, storage.ErrUserAlreadyExists) {
		log.Error("could not save user to database", sl.Err(err))
	}

	content := "📟 <b>Yo, welcome to Terminal Helper!</b>\n\n" +
		"This bot is developed to help you in @timetoterminal game.\n\n" +
		"<b>WARNING!</b> Take a notice, that this is not a hack or something like that. Bot just removes improper words, based on your attempts. All this stuff you can do manually.\n\n" +
		"The only thing, that can make your life a bit easier, words are sorted in such a way as to have the best chance of eliminating more words per attempt. So, its recommended to choose the <b>first (highest)</b> word in the list.\n\n" +
		"Although there are still games where you may not be able to guess a given word even after 4 attempts, but they're pretty rare"
	h.sendTextMessage(author.ID, content, GetMarkupNewGame())

	h.stages[author.ID] = None
}

func (h *Handler) CommandGame(u tgbotapi.Update) {
	author := u.Message.From

	game, exists := h.games[author.ID]
	if exists {
		content := "<b>You already have started game. Do you want to continue it?</b>\n\n<b>Words:</b>\n"
		for _, word := range game.AvailableWords() {
			content += fmt.Sprintf("<code>%s</code>\n", word)
		}
		h.sendTextMessage(author.ID, content, GetMarkupGameMenu())
	} else {
		h.sendTextMessage(author.ID, "Send me list of words in your $TERMINAL game", nil)
		h.stages[author.ID] = WaitingWordList
	}
}

func (h *Handler) CommandAdmin(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.CommandDataset"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("could not get user from database", sl.Err(err))
		return
	}

	if !user.IsAdmin {
		return
	}

	content := "<b>Admin Panel</b>\n\n"
	content += fmt.Sprintf("Logged in as <b>@%s</b>\n", author.UserName)
	content += fmt.Sprintf("<b>ID:</b> <code>%d</code>", author.ID)

	h.sendTextMessage(author.ID, content, GetMarkupAdmin())
}
