package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"terminal/internal/storage"
	"terminal/internal/terminal/dataset"
	"terminal/pkg/log/sl"
	"time"

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

	h.editMessage(author.ID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(game.AvailableWords()))
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
	h.editMessage(author.ID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(game.AvailableWords()))
}

func (h *Handler) CallbackDataset(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID
	log := h.log.With(
		slog.String("op", "handler.CallbackDataset"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("could not get user from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Something went wrong... Try again later</b>", nil)
		return
	}

	if !user.IsAdmin {
		h.editMessage(author.ID, messageID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	data, err := h.storage.GetDataset()
	if err != nil {
		log.Error("could not build dataset", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Failed to compose dataset</b>", nil)
		return
	}

	path, err := dataset.ExportDatasetToJSON(data)
	if err != nil {
		log.Error("could not export dataset to JSON", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Failed to compose dataset</b>", nil)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Error("could not open file", err)
		h.editMessage(author.ID, messageID, "<b>Failed to compose dataset</b>", nil)
		return
	}
	defer file.Close()

	reader := tgbotapi.FileReader{
		Name:   "0xterminal-dataset.json",
		Reader: file,
	}

	document := tgbotapi.NewDocument(author.ID, reader)
	document.ReplyMarkup = nil

	h.deleteMessage(author.ID, messageID)

	_, err = h.client.Send(document)
	if err != nil {
		log.Error("could not send dataset", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Failed to compose dataset</b>", nil)
	}

	log.Info("0xterminal dataset sent")

	content := "<b>Admin Panel</b>\n\n"
	content += fmt.Sprintf("Logged in as <b>@%s</b>\n", author.UserName)
	content += fmt.Sprintf("<b>ID:</b> <code>%d</code>", author.ID)

	h.sendTextMessage(author.ID, content, GetMarkupAdmin())

	os.Remove(path)
}

func (h *Handler) CallbackAdminPanel(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID
	log := h.log.With(
		slog.String("op", "handler.CallbackAdminPanel"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("could not get user from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Something went wrong... Try again later</b>", nil)
		return
	}

	if !user.IsAdmin {
		h.editMessage(author.ID, messageID, "<b>You are not permitted to use this action</b>", nil)
		return
	}

	content := "<b>Admin Panel</b>\n\n"
	content += fmt.Sprintf("Logged in as <b>@%s</b>\n", author.UserName)
	content += fmt.Sprintf("<b>ID:</b> <code>%d</code>", author.ID)

	h.editMessage(author.ID, messageID, content, GetMarkupAdmin())
}

func (h *Handler) CallbackStats(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID
	log := h.log.With(
		slog.String("op", "handler.CallbackStats"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("could not get user from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Something went wrong... Try again later</b>", GetMarkupBackToAdmin())
		return
	}

	if !user.IsAdmin {
		h.editMessage(author.ID, messageID, "<b>You are not permitted to use this action</b>", GetMarkupBackToAdmin())
		return
	}

	games, err := h.storage.GetAllGames()
	if err != nil {
		log.Error("could not get games from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Could not create statistics report</b>", GetMarkupBackToAdmin())
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("<b>All Time Statistics</b>\n\n<b>Total games:</b> %d\n", len(games)))

	gamesStats, err := h.storage.GetGamesToUserStatistics()
	if err != nil {
		log.Error("could not get games statistics from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Could not create statistics report</b>", GetMarkupBackToAdmin())
		return
	}

	for _, stat := range gamesStats {
		if stat.GamesPlayed != 0 {
			builder.WriteString(fmt.Sprintf(" - <b>%d</b> games played by @%s\n", stat.GamesPlayed, stat.Username))
		}
	}

	builder.WriteString("\n<b>Attempts ratio</b>\n")

	attemptCounts := make(map[int]int)

	for _, game := range games {
		attemptCounts[game.AttemptsAmount]++
	}

	attemptRatios := make(map[int]float64)
	for attempts, count := range attemptCounts {
		attemptRatios[attempts] = (float64(count) / float64(len(games))) * 100
	}

	attempts := make([]int, 0, len(attemptRatios))
	for attempt := range attemptRatios {
		attempts = append(attempts, attempt)
	}
	sort.Ints(attempts)

	for _, attempt := range attempts {
		if attemptCounts[attempt] == 1 {
			builder.WriteString(fmt.Sprintf(" - <b>%.2f%%</b> (%d game) completed in <b>%d</b> attempt", attemptRatios[attempt], attemptCounts[attempt], attempt))
		} else {
			builder.WriteString(fmt.Sprintf(" - <b>%.2f%%</b> (%d games) completed in <b>%d</b> attempt", attemptRatios[attempt], attemptCounts[attempt], attempt))
		}
		if attempt == 1 {
			builder.WriteString("\n")
		} else {
			builder.WriteString("s\n")
		}
	}

	usersAmount, err := h.storage.GetUsersCount()
	if err != nil {
		log.Error("could not get users amount from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Could not create statistics report</b>", GetMarkupBackToAdmin())
		return
	}

	builder.WriteString(fmt.Sprintf("\n<b>Total users:</b> %d", usersAmount))

	h.editMessage(author.ID, messageID, builder.String(), GetMarkupBackToAdmin())
}

func (h *Handler) CallbackDailyReport(u tgbotapi.Update) {
	author := u.CallbackQuery.From
	messageID := u.CallbackQuery.Message.MessageID
	log := h.log.With(
		slog.String("op", "handler.CallbackDailyReport"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	user, err := h.storage.GetUserByTelegramID(author.ID)
	if err != nil {
		log.Error("could not get user from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Something went wrong... Try again later</b>", GetMarkupBackToAdmin())
		return
	}

	if !user.IsAdmin {
		h.editMessage(author.ID, messageID, "<b>You are not permitted to use this action</b>", GetMarkupBackToAdmin())
		return
	}

	query := u.CallbackData()
	parts := strings.Split(query, ":")
	if len(parts) < 2 {
		log.Error("could not get date from callback query", slog.String("query", query))
		h.editMessage(author.ID, messageID, "<b>Something went wrong... Try again later</b>", GetMarkupBackToAdmin())
		return
	}

	var date time.Time
	if parts[1] == "today" {
		date = time.Now()
	} else {
		date, _ = time.Parse("02-01-2006", parts[1])
	}

	report, err := h.storage.GetDailyReport(date)
	if err != nil {
		log.Error("could not get daily report from database", sl.Err(err))
		h.editMessage(author.ID, messageID, "<b>Failed to get daily report</b>", GetMarkupBackToAdmin())
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

	content = fmt.Sprintf("<b>%s</b>\n\n<b>Games played:</b> %d\n", date.Format("2 January, 2006"), totalGames) + content

	content += fmt.Sprintf("<b>Joined users:</b> %d\n", len(report.JoinedUsers))
	for _, user := range report.JoinedUsers {
		content += fmt.Sprintf(" - @%s\n", user)
	}

	_, err = h.editMessage(author.ID, messageID, content, GetmarkupDailyReport(date))
	if err != nil {
		response := tgbotapi.NewCallback(u.CallbackQuery.ID, "No changes")
		h.client.Request(response)
	}
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

	game.SubmitAttempt(word, guessedLetters)

	if len(game.AvailableWords()) == 1 {
		delete(h.games, author.ID)
		h.editMessage(author.ID, messageID, fmt.Sprintf("<b>Target word:</b> <code>%s</code>", game.Target()), GetMarkupNewGame())

		// we'll assume that game is kinda spam, if initial words is less than 6
		if len(game.Words()) >= 6 {
			h.storage.SaveGame(author.ID, game.Words(), game.Target(), game.Attempts())
		}
		return
	}
	if len(game.AvailableWords()) == 0 {
		delete(h.games, author.ID)
		h.editMessage(author.ID, messageID, "<b>No matching words left.</b>\n\nTry again, may be you made a mistake?", nil)
		return
	}

	h.editMessage(author.ID, messageID, "<b>Pick one of the words in the list</b>", GetMarkupWords(h.games[author.ID].AvailableWords()))
}
