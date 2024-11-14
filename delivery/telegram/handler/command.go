package handler

import (
	"github.com/AsaHero/dastyor-bot/delivery/telegram/models"
	"gopkg.in/telebot.v3"
)

func (h *BotHandler) StartCommand(c telebot.Context) error {
	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(*models.ButtonRewrite),
	)

	return c.Send(models.MessageHomeDefault, keyboard)
}
