package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"terminal/internal/terminal"
	"terminal/pkg/log/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) PhotoMessage(u tgbotapi.Update) {
	author := u.Message.From
	log := h.log.With(
		slog.String("op", "handler.PhotoMessage"),
		slog.String("username", author.UserName),
		slog.String("id", strconv.FormatInt(author.ID, 10)),
	)

	// h.sendTextMessage(author.ID, "<b>To improve the quality of image recognition, please disable the effects in $TERMINAL\n\n</b>$TERMINAL -> settings -> effects -> turn off", nil)
	sticker, _ := h.sendSticker(author.ID, WaitingSticker)

	stage, exists := h.stages[author.ID]
	if !exists {
		h.stages[author.ID] = None
	}
	stage, _ = h.stages[author.ID]

	if stage == None {
		h.sendTextMessage(author.ID, "Use /newgame or click the button to start new $TERMINAL game", GetMarkupNewGame())
		return
	}

	photo := u.Message.Photo[len(u.Message.Photo)-1]

	fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
	file, err := h.client.GetFile(fileConfig)
	if err != nil {
		log.Error("failed to get file", sl.Err(err))
		h.sendTextMessage(u.Message.From.ID, "🚨 <b>Can't read words from this image</b>", nil)
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", h.client.Token, file.FilePath)

	if _, err = os.Stat("./.temp"); os.IsNotExist(err) {
		os.Mkdir("./.temp", 0755)
	}
	path := fmt.Sprintf("./.temp/%s.jpeg", photo.FileID)
	err = downloadFile(path, fileURL)
	if err != nil {
		log.Error("failed to download file", sl.Err(err))
		h.sendTextMessage(u.Message.From.ID, "🚨 <b>Can't read words from this image</b>", nil)
		return
	}

	words, err := h.ocr.ExtractWords(path)
	if err != nil {
		log.Error("can't read words from image", sl.Err(err))
		h.sendTextMessage(u.Message.From.ID, "🚨 <b>Can't read words from this image</b>", nil)
		return
	}

	err = os.Remove(path)
	if err != nil {
		h.log.Error("failed to delete temporary file", sl.Err(err))
	}

	h.deleteMessage(author.ID, sticker.MessageID)

	if len(words) < 6 {
		var content string
		if len(words) == 0 {
			content = "<b>No words recognized</b>"
		} else {
			content = "<b>Recognized words:</b>\n\n"
			for _, word := range words {
				content += fmt.Sprintf("<code>%s</code>\n", word)
			}
		}
		h.sendTextMessage(author.ID, content, nil)
		h.sendTextMessage(author.ID, "<b>According to the $TERMINAL rules, the word list must consist of at least 6 words</b>\n\nSend me list of words in your $TERMINAL game", nil)
		return
	}

	game, err := terminal.New(words)
	if errors.Is(err, terminal.ErrDifferentWordsLength) {
		content := "<b>Recognized words:</b>\n\n"
		for _, word := range words {
			content += fmt.Sprintf("<code>%s</code>\n", word)
		}
		h.sendTextMessage(author.ID, content, nil)
		h.sendTextMessage(author.ID, "<b>According to the $TERMINAL rules, the word list should only consist of words of the same length</b>\n\nSend me list of words in your $TERMINAL game", nil)
		return
	}
	h.games[author.ID] = game
	h.stages[author.ID] = None

	answer, err := h.storage.TryFindAnswer(words)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("game with this word list not found")
		} else {
			log.Error("could not get answer from database", sl.Err(err))
		}
	}
	if answer != "" {
		h.sendTextMessage(author.ID, "<b>Found game with similar words list</b>\n\nProbably, the target is <code>"+answer+"</code>", nil)
	}

	h.sendTextMessage(author.ID, fmt.Sprintf("<b>Pick one of %d words in the list</b>", len(words)), GetMarkupWords(h.games[author.ID].AvailableWords))
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
