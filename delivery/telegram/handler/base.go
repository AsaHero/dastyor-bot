package handler

import (
	"context"

	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"github.com/AsaHero/dastyor-bot/internal/usecase/llm"
	"github.com/AsaHero/dastyor-bot/internal/usecase/sessions"
	"github.com/AsaHero/dastyor-bot/internal/usecase/users"
	"github.com/AsaHero/dastyor-bot/pkg/config"
	telegram_bot "github.com/AsaHero/dastyor-bot/pkg/telegram-bot"
	"gopkg.in/telebot.v3"
)

type BotHandler struct {
	cfg             *config.Config
	bot             *telegram_bot.TelegramBot
	userService     users.Users
	sessionsService sessions.Sessions
	llmService      llm.LLM
}

type Options struct {
	UserService     users.Users
	SessionsService sessions.Sessions
	LlmService      llm.LLM
}

func NewBotHandler(cfg *config.Config, bot *telegram_bot.TelegramBot, options *Options) *BotHandler {
	return &BotHandler{
		cfg:             cfg,
		bot:             bot,
		userService:     options.UserService,
		sessionsService: options.SessionsService,
		llmService:      options.LlmService,
	}
}

func (h *BotHandler) SessionsAssigner(c telebot.Context, sessionsType telegram_bot.BotSessionType, sessionsState telegram_bot.BotState) error {
	user := c.Sender()
	if user != nil {
		err := h.sessionsService.Set(context.Background(), &entity.Sessions{
			UserID: user.ID,
			Type:   sessionsType,
			State:  sessionsState,
		})
		if err != nil {
			return inerr.WithMessage(err, "failed to assign userID: %d", user.ID)
		}
	}

	return nil
}
