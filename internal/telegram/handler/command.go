package handler

import (
	"fmt"
	"os"
	"terminal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	userID := u.Message.From.ID
	h.sendTextMessage(userID, "<b>Yoo</b>\nUse /newgame or click the button to start new $TERMINAL game", GetMarkupNewGame())
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
