package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"terminal/internal/storage"
	"terminal/internal/terminal/dataset"
	"terminal/pkg/git"
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

	// TODO(#27): Add build info to /a command such as last commit, version etc.
	content := "<b>Admin Panel</b>\n\n"
	content += fmt.Sprintf("Logged in as <b>@%s</b>\n", author.UserName)
	content += fmt.Sprintf("<b>ID:</b> <code>%d</code>\n\n", author.ID)
	content += fmt.Sprintf("<b>Build:</b>\n")
	content += fmt.Sprintf("Commit: <a href=\"https://github.com/jus1d/0xterminal-helper/tree/%s\">%s</a>\n", git.LatestCommit(), git.LatestShortenedCommit())
	content += fmt.Sprintf("Branch: <code>%s</code>", git.CurrentBranch())

	h.sendTextMessage(author.ID, content, GetMarkupAdmin())
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
		log.Error("could not get user from database", sl.Err(err))
		h.sendTextMessage(author.ID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	if !user.IsAdmin {
		h.sendTextMessage(author.ID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	data, err := h.storage.GetDataset()
	if err != nil {
		log.Error("could not build dataset", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to compose dataset</b>", nil)
		return
	}

	path, err := dataset.ExportDatasetToJSON(data)
	if err != nil {
		log.Error("could not export dataset to JSON", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to compose dataset</b>", nil)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Error("could not open file", err)
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
		log.Error("could not send dataset", sl.Err(err))
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
		log.Error("could not get user from database", sl.Err(err))
		h.sendTextMessage(author.ID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	if !user.IsAdmin {
		h.sendTextMessage(author.ID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	report, err := h.storage.GetDailyReport()
	if err != nil {
		log.Error("could not get daily report from database", sl.Err(err))
		h.sendTextMessage(author.ID, "🚨 <b>Failed to get daily report</b>", nil)
		return
	}

	var content string

	totalGames := 0
	for i, stat := range report.Stats {
		if stat.GamesPlayed == 1 {
			content += fmt.Sprintf(" - <b>%d</b> game by @%s\n", stat.GamesPlayed, stat.Username)
		} else {
			content += fmt.Sprintf(" - <b>%d</b> games by @%s\n", stat.GamesPlayed, stat.Username)
		}
		totalGames += stat.GamesPlayed
		if i == len(report.Stats)-1 {
			content += "\n"
		}
	}

	content = fmt.Sprintf("<b>%s</b>\n\n<b>Games played:</b> %d\n", time.Now().Format("2 January, 2006"), totalGames) + content

	content += fmt.Sprintf("<b>Joined users:</b> %d\n", len(report.JoinedUsers))
	for _, user := range report.JoinedUsers {
		content += fmt.Sprintf(" - @%s\n", user)
	}

	h.sendTextMessage(author.ID, content, nil)
}
