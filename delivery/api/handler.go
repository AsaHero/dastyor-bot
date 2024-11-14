package api

import (
	"net/http"

	"github.com/AsaHero/dastyor-bot/pkg/config"
	telegram_bot "github.com/AsaHero/dastyor-bot/pkg/telegram-bot"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg         *config.Config
	telegramBot *telegram_bot.TelegramBot
}

type HandlerOptions struct {
	Config      *config.Config
	TelegramBot *telegram_bot.TelegramBot
}

func NewHandler(g *gin.RouterGroup, options HandlerOptions) {
	handler := &Handler{
		cfg:         options.Config,
		telegramBot: options.TelegramBot,
	}

	g.POST("/telegram", gin.WrapF(handler.telegramWebhook))
}

func (r *Handler) telegramWebhook(write http.ResponseWriter, req *http.Request) {
	r.telegramBot.HandleUpdate(write, req)
}
