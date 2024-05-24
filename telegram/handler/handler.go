package handler

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

type Handler struct {
	client      *tgbotapi.BotAPI
	isDebugMode bool
	games       map[int64]*terminal.Game
	stages      map[int64]Stage
}

func New(client *tgbotapi.BotAPI, isDebugMode bool) *Handler {
	return &Handler{
		client:      client,
		isDebugMode: isDebugMode,
		games:       make(map[int64]*terminal.Game, 0),
		stages:      make(map[int64]Stage, 0),
	}
}

func (h *Handler) CommandStart(u tgbotapi.Update) {
	userID := u.Message.From.ID
	h.sendTextMessage(userID, "Use /game to register a game")
	h.stages[userID] = None
}

func (h *Handler) CommandGame(u tgbotapi.Update) {
	userID := u.Message.From.ID
	h.sendTextMessage(userID, "Send me list of words in your $TERMINAL game")
	h.stages[userID] = WaitingWordList
}

func (h *Handler) TextMessage(u tgbotapi.Update) {
	userID := u.Message.From.ID
	stage, exists := h.stages[userID]
	if exists && stage == WaitingWordList {
		words := strings.Split(u.Message.Text, "\n")
		game, err := terminal.New(words)
		if errors.Is(err, terminal.ErrDifferentWordsLength) {
			h.sendTextMessage(userID, "Words must be the same length")
			return
		}
		h.games[userID] = game
		h.stages[userID] = WaitingAttempt

		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(h.games[userID].Words))
		for i, word := range h.games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		h.sendTextMessage(userID, content)
	} else if exists && stage == WaitingAttempt {
		parts := strings.Split(u.Message.Text, " ")
		if len(parts) < 2 {
			h.sendTextMessage(userID, "U invalid. Use: word guessed-letters-amount")
			return
		}
		word := parts[0]
		guessedLetters, err := strconv.Atoi(parts[1])
		if err != nil {
			h.sendTextMessage(userID, "Guessed letters amount should be integer. Use: word guessed-letters-amount")
			return
		}
		attempt := terminal.Attempt{
			Word:           word,
			GuessedLetters: guessedLetters,
		}
		h.games[userID].CommitAttempt(attempt)
		content := fmt.Sprintf("<b>Available %d words:</b>\n", len(h.games[userID].Words))
		for i, word := range h.games[userID].Words {
			content += fmt.Sprintf("#%d: <code>%s</code>\n", i+1, word)
		}
		h.sendTextMessage(userID, content)
		if len(h.games[userID].Words) <= 1 {
			h.stages[userID] = None
		}
	}
}

func (h *Handler) sendTextMessage(chatID int64, content string) {
	message := tgbotapi.NewMessage(chatID, content)
	message.ParseMode = tgbotapi.ModeHTML

	_, err := h.client.Send(message)
	if err != nil {
		log.Printf("ERROR: could not send message to ID: %d, error: %s\n", chatID, err.Error())
		return
	}
	if h.isDebugMode {
		log.Printf("DEBUG: message sent to ID: %d, content: %s\n", chatID, content)
	}
}
