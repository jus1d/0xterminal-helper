package handler

import (
	"fmt"
	"os"
	"terminal/pkg/log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CommandStart(u tgbotapi.Update) {
	user := u.Message.From

	_, err := h.storage.CreateUser(user.ID, user.UserName, user.FirstName, user.LastName)
	if err != nil {
		log.Error("could not create new user", err, log.WithString("username", user.UserName))
	}

	content := "📟 <b>Yo, welcome to Terminal Helper!</b>\n\n" +
		"This bot is developed to help you in @timetoterminal game.\n\n" +
		"<b>WARNING!</b> Take a notice, that this is not a hack or something like that. Bot just removes improper words, based on your attempts. All this stuff you can do manually.\n\n" +
		"The only thing, that can make your life a bit easier, words are sorted in such a way as to have the best chance of eliminating more words per attempt. So, its recommended to choose the <b>first (highest)</b> word in the list.\n\n" +
		"Although there are still games where you may not be able to guess a given word even after 4 attempts, but they're pretty rare"
	h.sendTextMessage(user.ID, content, GetMarkupNewGame())

	h.stages[user.ID] = None
}

func (h *Handler) CommandGame(u tgbotapi.Update) {
	userID := u.Message.From.ID

	_, err := h.storage.GetUserByTelegramID(userID)
	if err != nil {
		log.Error("could not get user from database", err, log.WithInt64("telegram_id", userID))
		h.sendTextMessage(userID, "<b>It seems that you are new here</b>\n\nUse /start to start the bot", nil)
		return
	}

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

	user, err := h.storage.GetUserByTelegramID(userID)
	if err != nil {
		log.Error("could not get user from database", err, log.WithInt64("telegram_id", userID))
		h.sendTextMessage(userID, "<b>It seems that you are new here</b>\n\nUse /start to start the bot", nil)
		return
	}

	if !user.IsAdmin {
		h.sendTextMessage(userID, "<b>You are not permitted to use this command</b>", nil)
		return
	}

	// TODO: remove creating JSON file from storage to other package
	path := time.Now().Format("./dataset-02-01-2006.json")
	err = h.storage.ParseGamesToJsonFile(path)
	if err != nil {
		h.sendTextMessage(userID, "🚨 <b>Failed to compose dataset</b>", nil)
		log.Error("failed to compose dataset", err)
		return
	}

	file, err := os.Open(path)
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

	os.Remove(path)
}
