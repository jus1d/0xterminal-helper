package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"terminal/internal/terminal/dataset"
	"terminal/pkg/log/sl"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.CommandStart"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	_, err := h.storage.GetUserByTelegramID(author.ID)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := h.storage.CreateUser(author.ID, author.UserName, author.FirstName, author.LastName)
		if err != nil {
			log.Error("failed to create new user in database", sl.Err(err))
		}
	} else if err != nil {
		log.Error("failed to get user from database", sl.Err(err))
		return
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
	log := h.log.With(
		slog.String("op", "handler.CommandGame"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	_, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("failed to get user from database", sl.Err(err))
		h.sendTextMessage(author.ID, "<b>It seems that you are new here</b>\n\nUse /start to start the bot", nil)
		return
	}

	game, exists := h.games[author.ID]
	if exists {
		content := "<b>You already have started game. Do you want to continue it?</b>\n\n<b>Words:</b>\n"
		for _, word := range game.AvailableWords {
			content += fmt.Sprintf("<code>%s</code>\n", word)
		}
		h.sendTextMessage(author.ID, content, GetMarkupGameMenu())
	} else {
		h.sendTextMessage(author.ID, "Send me list of words in your $TERMINAL game", nil)
		h.stages[author.ID] = WaitingWordList
	}
}

func (h *Handler) CommandDataset(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.CommandDataset"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("failed to get user from database", sl.Err(err))
		h.sendTextMessage(author.ID, "<b>It seems that you are new here</b>\n\nUse /start to start the bot", nil)
		return
	}

	if !user.IsAdmin {
		h.sendTextMessage(author.ID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	data, err := h.storage.GetDataset()
	if err != nil {
		log.Error("failed to build dataset", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to compose dataset</b>", nil)
		return
	}

	path, err := dataset.ExportDatasetToJSON(data)
	if err != nil {
		log.Error("failed to export dataset to JSON", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to compose dataset</b>", nil)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Error("failed to open file", err)
		h.sendTextMessage(author.ID, "🚨 <b>Failed to compose dataset</b>", nil)
		return
	}
	defer file.Close()

	reader := tgbotapi.FileReader{
		Name:   "0xterminal-dataset.json",
		Reader: file,
	}

	document := tgbotapi.NewDocument(author.ID, reader)

	_, err = h.client.Send(document)
	if err != nil {
		log.Error("failed to send dataset", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to compose dataset</b>", nil)
	}

	log.Info("0xterminal dataset sent")

	os.Remove(path)
}

func (h *Handler) CommandDailyReport(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.CommandDailyReport"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("failed to get user from database", sl.Err(err))
		h.sendTextMessage(author.ID, "<b>It seems that you are new here</b>\n\nUse /start to start the bot", nil)
		return
	}

	if !user.IsAdmin {
		h.sendTextMessage(author.ID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	report, err := h.storage.GetDailyReport()
	if err != nil {
		log.Error("failed tp get daily report from database", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to get daily report</b>", nil)
		return
	}

	var content string

	totalGames := 0
	for i, stat := range report.Stats {
		content += fmt.Sprintf(" - <b>%d</b> by @%s\n", stat.GamesPlayed, stat.Username)
		totalGames += stat.GamesPlayed
		if i == len(report.Stats)-1 {
			content += "\n"
		}
	}

	content = fmt.Sprintf("<b>%s</b>\n\n<b>Games played: %d</b>\n", time.Now().Format("2 January, 2006"), totalGames) + content

	content += fmt.Sprintf("<b>Joined users:</b> %d\n", len(report.JoinedUsers))
	for _, user := range report.JoinedUsers {
		content += fmt.Sprintf(" - @%s\n", user)
	}

	h.sendTextMessage(author.ID, content, nil)
}
