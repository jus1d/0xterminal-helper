package handler

import (
	"fmt"
	"os"
	"terminal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	userID := u.Message.From.ID

	content := "📟 <b>Yo, welcome to Terminal Helper!</b>\n\n" +
		"This bot is developed to help you in @timetoterminal game.\n\n" +
		"<b>WARNING!</b> Take a notice, that this is not a hack or something like that. Bot just removes improper words, based on your attempts. All this stuff you can do manually.\n\n" +
		"The only thing, that can make your life a bit easier, words are sorted in such a way as to have the best chance of eliminating more words per attempt. So, its recommended to choose the <b>first (highest)</b> word in the list.\n\n" +
		"Although there are still games where you may not be able to guess a given word even after 4 attempts, but they're pretty rare"
	h.sendTextMessage(userID, content, GetMarkupNewGame())

	h.stages[userID] = None
}

func (h *Handler) CommandGame(u tgbotapi.Update) {
	userID := u.Message.From.ID
	game, exists := h.games[userID]
	if exists {
		content := "<b>You already have started game. Do you want to continue it?</b>\n\n<b>Words:</b>\n"
		for _, word := range game.AvailableWords {
			content += fmt.Sprintf("<code>%s</code>\n", word)
		}
		h.sendTextMessage(userID, content, GetMarkupGameMenu())
	} else {
		h.sendTextMessage(userID, "Send me list of words in your $TERMINAL game", nil)
		h.stages[userID] = WaitingWordList
	}
}

func (h *Handler) CommandDataset(u tgbotapi.Update) {
	userID := u.Message.From.ID

	file, err := os.Open("./storage/data.json")
	if err != nil {
		log.Error("could not open file", err)
		h.sendTextMessage(userID, "Could not send file", nil)
	}
	defer file.Close()

	reader := tgbotapi.FileReader{
		Name:   "0xterminal-dataset.json",
		Reader: file,
	}

	document := tgbotapi.NewDocument(userID, reader)

	_, err = h.client.Send(document)
	if err != nil {
		log.Info("dataset sent", log.WithInt64("id", userID), log.WithString("username", u.Message.From.UserName))
	}
}
