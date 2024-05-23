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
	IsDebugMode bool
	Games       map[int64]*terminal.Game
	Stages      map[int64]Stage
	*tgbotapi.BotAPI
}

func New(token string, isDebugMode bool) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("ERROR: could not start the bot")
	}

	return &Bot{
		IsDebugMode: isDebugMode,
		Games:       make(map[int64]*terminal.Game, 0),
		Stages:      make(map[int64]Stage, 0),
		BotAPI:      bot,
	}
}

func (bot *Bot) Run() {
	log.Printf("Authorized on account @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		bot.handleUpdate(update)
	}
}

func (bot *Bot) handleUpdate(u tgbotapi.Update) {
	if u.Message != nil {
		log.Printf("[@%s] says: `%s`", u.Message.From.UserName, u.Message.Text)

		switch u.Message.Text {
		case "/start":
			bot.handleStartCommand(u)
		case "/game":
			bot.handleGameCommand(u)
		default:
			bot.handleTextMessage(u)
		}
	}
}

func (bot *Bot) handleStartCommand(u tgbotapi.Update) {
	userID := u.Message.From.ID
	bot.sendMessage(userID, "Hello, Im gonna help you with your TERMINAL games. Use /game to start")
	bot.Stages[userID] = None
}

func (bot *Bot) handleGameCommand(u tgbotapi.Update) {
	userID := u.Message.From.ID
	_, exists := bot.Games[userID]
	if exists {
		// TODO: Add buttons: continue | start new game
		bot.sendMessage(userID, "[TODO] You already have started game. Do you want to continue?")
		return
	}

	bot.sendMessage(userID, "Send me list of words in your TERMINAL game")
	bot.Stages[userID] = WaitingWordList
}

func (bot *Bot) handleTextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID
	stage, exists := bot.Stages[userID]
	if exists && stage == WaitingWordList {
		words := strings.Split(u.Message.Text, "\n")
		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			bot.sendMessage(userID, "Words should not be different length")
			return
		}
		bot.Games[userID] = game
		bot.Stages[userID] = WaitingAttempt

		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(bot.Games[userID].Words))
		for i, word := range bot.Games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		bot.sendMessage(userID, content)
	} else if exists && stage == WaitingAttempt {
		parts := strings.Split(u.Message.Text, " ")
		word := parts[0]
		guessedLetters, _ := strconv.Atoi(parts[1])
		attempt := terminal.Attempt{
			Word:           word,
			GuessedLetters: guessedLetters,
		}
		bot.Games[userID].Attempts = append(bot.Games[userID].Attempts, &attempt)
		bot.Games[userID].UpdateWords()
		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(bot.Games[userID].Words))
		for i, word := range bot.Games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		bot.sendMessage(userID, content)
		if len(bot.Games[userID].Words) == 1 {
			bot.Stages[userID] = None
		}
	}
}

func (bot *Bot) sendMessage(chatID int64, content string) {
	message := tgbotapi.NewMessage(chatID, content)
	message.ParseMode = tgbotapi.ModeHTML

	_, err := bot.Send(message)
	if err != nil {
		log.Printf("ERROR: could not send message to ID: %d, error: %s\n", chatID, err.Error())
		return
	}
	if bot.IsDebugMode {
		log.Printf("DEBUG: message sent to ID: %d, content: %s\n", chatID, content)
	}
}
