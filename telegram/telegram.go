package telegram

import (
	"log"
	"terminal/telegram/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	client  *tgbotapi.BotAPI
	handler *handler.Handler
}

func New(token string, isDebugMode bool) *Bot {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("ERROR: could not start the bot")
	}

	return &Bot{
		client:  client,
		handler: handler.New(client, isDebugMode),
	}
}

func (b *Bot) Run() {
	log.Printf("INFO: authorized on account @%s", b.client.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.client.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	if u.Message != nil {
		log.Printf("[@%s] says: `%s`", u.Message.From.UserName, u.Message.Text)

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
		switch query {
		case "game-continue":
			b.handler.CallbackContinueGame(u)
		case "start-new-game":
			b.handler.CallbackStartNewGame(u)
		}
	}
}
