package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *Handler) CallbackContinueGame(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	h.stages[userID] = WaitingAttempt
	h.editMessage(userID, u.CallbackQuery.Message.MessageID, "Send me your attempt", nil)
}

func (h *Handler) CallbackStartNewGame(u tgbotapi.Update) {
	userID := u.CallbackQuery.From.ID
	h.stages[userID] = WaitingWordList
	h.editMessage(userID, u.CallbackQuery.Message.MessageID, "Send me list of words in your $TERMINAL game", nil)
}
