package telegram

import (
	"context"

	"github.com/AsaHero/dastyor-bot/delivery/telegram/handler"
	"github.com/AsaHero/dastyor-bot/delivery/telegram/models"
	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"github.com/AsaHero/dastyor-bot/internal/usecase/sessions"
	"github.com/AsaHero/dastyor-bot/pkg/config"
	telegram_bot "github.com/AsaHero/dastyor-bot/pkg/telegram-bot"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type BotRouter struct {
	cfg             *config.Config
	bot             *telegram_bot.TelegramBot
	handler         *handler.BotHandler
	sessionsService sessions.Sessions
}

func NewRouter(cfg *config.Config, tgBot *telegram_bot.TelegramBot, options *handler.Options) *BotRouter {
	tgBot.Use(middleware.Recover())
	tgBot.Use(middleware.AutoRespond())
	tgBot.Use(middleware.Logger())

	r := &BotRouter{
		cfg:             cfg,
		bot:             tgBot,
		handler:         handler.NewBotHandler(cfg, tgBot, options),
		sessionsService: options.SessionsService,
	}

	tgBot.Handle("/start", r.handler.StartCommand)
	tgBot.Handle(models.ButtonHome, r.handler.StartCommand)
	tgBot.Handle(telebot.OnText, r.TextMessageHandler)
	tgBot.Handle(models.ButtonRewrite, r.handler.RewriteInit, r.SessionsAssignerMiddleware(entity.REWRITE, entity.REWRITE_waiting))

	// Rewrite session handler
	tgBot.Subscribe(entity.REWRITE, entity.REWRITE_waiting, r.handler.RewriteHandle, r.SessionsAssignerMiddleware(entity.REWRITE, entity.REWRITE_answering))

	return r
}

func (r BotRouter) TextMessageHandler(c telebot.Context) error {
	ctx := context.Background()
	user := c.Sender()
	if user == nil {
		return nil
	}

	session, err := r.sessionsService.Get(ctx, user.ID)
	if err != nil {
		return inerr.Err(err)
	}

	r.bot.Publish(session.Type, session.State, c)

	return nil
}

func (h *BotRouter) SessionsAssignerMiddleware(sessionsType telegram_bot.BotSessionType, sessionsState telegram_bot.BotState) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			user := ctx.Sender()
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

			return next(ctx)
		}
	}
}
