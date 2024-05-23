package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"terminal/terminal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stage uint8

const (
	None = iota
	WaitingWordList
	WaitingAttempt
)

type Bot struct {
	isDebugMode bool
	games       map[int64]*terminal.Game
	stages      map[int64]Stage
	bot         *tgbotapi.BotAPI
}

func New(token string, isDebugMode bool) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("ERROR: could not start the bot")
	}

	return &Bot{
		isDebugMode: isDebugMode,
		games:       make(map[int64]*terminal.Game, 0),
		stages:      make(map[int64]Stage, 0),
		bot:         bot,
	}
}

func (b *Bot) Run() {
	log.Printf("INFO: authorized on account @%s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	if u.Message != nil {
		log.Printf("[@%s] says: `%s`", u.Message.From.UserName, u.Message.Text)

		switch u.Message.Text {
		case "/start":
			b.handleStartCommand(u)
		case "/game":
			b.handleGameCommand(u)
		default:
			b.handleTextMessage(u)
		}
	}
}

func (b *Bot) handleStartCommand(u tgbotapi.Update) {
	userID := u.Message.From.ID
	b.sendMessage(userID, "Use /game to register a game")
	b.stages[userID] = None
}

func (b *Bot) handleGameCommand(u tgbotapi.Update) {
	userID := u.Message.From.ID

	b.sendMessage(userID, "Send me list of words in your $TERMINAL game")
	b.stages[userID] = WaitingWordList
}

func (b *Bot) handleTextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID
	stage, exists := b.stages[userID]
	if exists && stage == WaitingWordList {
		words := strings.Split(u.Message.Text, "\n")
		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			b.sendMessage(userID, "Words must be the same length")
			return
		}
		b.games[userID] = game
		b.stages[userID] = WaitingAttempt

		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(b.games[userID].Words))
		for i, word := range b.games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		b.sendMessage(userID, content)
	} else if exists && stage == WaitingAttempt {
		parts := strings.Split(u.Message.Text, " ")
		if len(parts) < 2 {
			b.sendMessage(userID, "U invalid. Use: word guessed-letters-amount")
			return
		}
		word := parts[0]
		guessedLetters, _ := strconv.Atoi(parts[1])
		attempt := terminal.Attempt{
			Word:           word,
			GuessedLetters: guessedLetters,
		}
		b.games[userID].CommitAttempt(attempt)
		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(b.games[userID].Words))
		for i, word := range b.games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		b.sendMessage(userID, content)
		if len(b.games[userID].Words) == 1 {
			b.stages[userID] = None
		}
	}
}

func (b *Bot) sendMessage(chatID int64, content string) {
	message := tgbotapi.NewMessage(chatID, content)
	message.ParseMode = tgbotapi.ModeHTML

	_, err := b.bot.Send(message)
	if err != nil {
		log.Printf("ERROR: could not send message to ID: %d, error: %s\n", chatID, err.Error())
		return
	}
	if b.isDebugMode {
		log.Printf("DEBUG: message sent to ID: %d, content: %s\n", chatID, content)
	}
}
