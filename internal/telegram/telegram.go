package telegram

import (
	"log/slog"
	"os"
	"strings"
	"terminal/internal/config"
	"terminal/internal/ocr"
	"terminal/internal/storage"
	"terminal/internal/telegram/handler"
	"terminal/pkg/log/sl"
	"terminal/pkg/str"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	log     *slog.Logger
	client  *tgbotapi.BotAPI
	handler *handler.Handler
}

func New(log *slog.Logger, conf config.Telegram, st storage.Storage, o *ocr.Client) *Bot {
	client, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Error("failed to start the bot", sl.Err(err))
		os.Exit(1)
	}

	return &Bot{
		log:     log,
		client:  client,
		handler: handler.New(log, client, st, o),
	}
}

func (b *Bot) Run() {
	b.log.Info("bot authorized into telegram API", slog.String("username", b.client.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.client.GetUpdatesChan(u)

	for update := range updates {
		go b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	log := b.log.With(
		slog.String("op", "bot.handleUpdate"),
	)

	if u.Message != nil {
		if u.Message.Photo != nil {
			log.Info("photo message received", slog.Int64("id", u.Message.From.ID), slog.String("username", u.Message.From.UserName))

			b.handler.PhotoMessage(u)
			return
		}

		log.Info("text message received", slog.String("content", str.Unescape(u.Message.Text)), slog.Int64("id", u.Message.From.ID), slog.String("username", u.Message.From.UserName))

		switch u.Message.Text {
		case "/start":
			b.handler.CommandStart(u)
		case "/newgame":
			b.handler.CommandGame(u)
		case "/dataset":
			b.handler.CommandDataset(u)
		case "/dailyreport":
			b.handler.CommandDailyReport(u)
		case "/a":
			b.handler.CommandAdmin(u)
		default:
			b.handler.TextMessage(u)
		}

	}
	if u.CallbackQuery != nil {
		query := u.CallbackData()
		log.Info("callback received", slog.String("query", query), slog.Int64("id", u.CallbackQuery.From.ID), slog.String("username", u.CallbackQuery.From.UserName))

		switch {
		case query == "game-continue":
			b.handler.CallbackContinueGame(u)
		case query == "start-new-game":
			b.handler.CallbackStartNewGame(u)
		case query == "words-list":
			b.handler.CallbackWordsList(u)
		case strings.HasPrefix(query, "choose-word:"):
			b.handler.CallbackChooseWord(u)
		case strings.HasPrefix(query, "choose-guessed-letters:"):
			b.handler.CallbackChooseGuessedLetters(u)
		}
	}
}
