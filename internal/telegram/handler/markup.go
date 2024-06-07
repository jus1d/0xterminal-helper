package handler

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMarkupAdmin() *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Daily Report", time.Now().Format("daily-report:02-01-2006")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Dataset", "dataset"),
		),
	)
	return &markup
}

func GetmarkupDailyReport(date time.Time) *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("«", date.Add(-24*time.Hour).Format("daily-report:02-01-2006")),
			tgbotapi.NewInlineKeyboardButtonData("↻", date.Format("daily-report:02-01-2006")),
			tgbotapi.NewInlineKeyboardButtonData("»", date.Add(24*time.Hour).Format("daily-report:02-01-2006")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Today", time.Now().Format("daily-report:02-01-2006")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Back", "admin-panel"),
		),
	)
	return &markup
}

func GetMarkupBackToAdmin() *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Back", "admin-panel"),
		),
	)
	return &markup
}

func GetMarkupGameMenu() *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Continue", "game-continue"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Start new game", "start-new-game"),
		),
	)
	return &markup
}

func GetMarkupWords(words []string) *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, word := range words {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(word, fmt.Sprintf("choose-word:%s", word)),
		))
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup
}

func GetMarkupGuessedLetters(word string) *tgbotapi.InlineKeyboardMarkup {
	if len(word) > 21 {
		return getHugeMarkupGuessedLetters(word)
	}
	// rules describes how many buttons should be in each row depending on the word length
	rules := map[int][]int{
		1:  {1},
		2:  {2},
		3:  {3},
		4:  {4},
		5:  {3, 2},
		6:  {3, 3},
		7:  {4, 3},
		8:  {4, 4},
		9:  {5, 4},
		10: {5, 5},
		11: {6, 5},
		12: {6, 6},
		13: {5, 5, 3},
		14: {5, 5, 4},
		15: {5, 5, 5},
		16: {6, 5, 5},
		17: {5, 5, 4, 3},
		18: {5, 5, 5, 4},
		19: {5, 5, 5, 4},
		20: {5, 5, 5, 5},
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("None", fmt.Sprintf("choose-guessed-letters:%s:%d", word, 0))})

	rule := rules[len(word)-1]
	i := 1
	for ri := range rule {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		for rule[ri] != 0 && i < len(word) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i), fmt.Sprintf("choose-guessed-letters:%s:%d", word, i)))
			i++
			rule[ri]--
		}
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("All", fmt.Sprintf("choose-guessed-letters:%s:%d", word, len(word)))})

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Back", "words-list"),
	))

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup
}

func getHugeMarkupGuessedLetters(word string) *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("None", fmt.Sprintf("choose-guessed-letters:%s:%d", word, 0))})

	for i := 1; i < len(word); i += 5 {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		for j := 0; len(row) < 5 && i+j < len(word); j++ {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+j), fmt.Sprintf("choose-guessed-letters:%s:%d", word, i+j)))
		}
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("All", fmt.Sprintf("choose-guessed-letters:%s:%d", word, len(word)))})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup
}

func GetMarkupNewGame() *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Start new game", "start-new-game"),
		),
	)
	return &markup
}
