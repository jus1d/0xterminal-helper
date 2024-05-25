package telegram

import (
	"strings"
	"terminal/internal/telegram/handler"
	"terminal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	client  *tgbotapi.BotAPI
	handler *handler.Handler
}

func New(token string) *Bot {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("could not start the bot", err)
	}

	return &Bot{
		client:  client,
		handler: handler.New(client),
	}
}

func (b *Bot) Run() {
	log.Info("bot authorized into telegram API", log.WithString("account", b.client.Self.FirstName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.client.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	if u.Message != nil {
		log.Info("message recieved", log.WithString("username", u.Message.From.UserName), log.WithInt64("id", u.Message.From.ID), log.WithString("content", u.Message.Text))

		switch u.Message.Text {
		case "/start":
			b.handler.CommandStart(u)
		case "/game":
			b.handler.CommandGame(u)
		default:
			b.handler.TextMessage(u)
		}
	}
	if u.CallbackQuery != nil {
		query := u.CallbackData()
		log.Info("callback recieved", log.WithString("username", u.CallbackQuery.From.UserName), log.WithInt64("id", u.CallbackQuery.From.ID), log.WithString("query", query))

		switch {
		case query == "game-continue":
			b.handler.CallbackContinueGame(u)
		case query == "start-new-game":
			b.handler.CallbackStartNewGame(u)
		case strings.HasPrefix(query, "choose-word:"):
			b.handler.CallbackChooseWord(u)
		case strings.HasPrefix(query, "choose-guessed-letters:"):
			b.handler.CallbackChooseGuessedLetters(u)
		}
	}
}
