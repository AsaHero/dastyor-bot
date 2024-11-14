package telegram_bot

import (
	"log"
	"net/http"

	"gopkg.in/telebot.v3"
)

type BotState string
type BotSessionType string

type TelegramBot struct {
	*telebot.Bot
	webhook         *telebot.Webhook
	sessionHandlers map[BotSessionType]map[BotState]HandlerConfig
}

type HandlerConfig struct {
	handler     func(telebot.Context) error
	middlewares []telebot.MiddlewareFunc
}

func NewTelegramBot(token string, webhookURL string) *TelegramBot {
	webhook := &telebot.Webhook{
		AllowedUpdates: []string{
			"callback_query",
			"edited_message",
			"message",
			"pre_checkout_query",
		},
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: webhookURL,
		},
		DropUpdates: true,
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:       token,
		Poller:      webhook,
		Synchronous: false,
		Verbose:     false,
		ParseMode:   telebot.ModeMarkdown,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &TelegramBot{
		Bot:             bot,
		webhook:         webhook,
		sessionHandlers: make(map[BotSessionType]map[BotState]HandlerConfig),
	}
}

func (b *TelegramBot) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	b.webhook.ServeHTTP(w, r)
}

// Subscribe adds a handler with optional middlewares for a specific session type and state
func (b *TelegramBot) Subscribe(sessionType BotSessionType, sessionState BotState, handler func(telebot.Context) error, middlewares ...telebot.MiddlewareFunc) {
	if _, ok := b.sessionHandlers[sessionType]; !ok {
		b.sessionHandlers[sessionType] = make(map[BotState]HandlerConfig)
	}

	if _, ok := b.sessionHandlers[sessionType][sessionState]; !ok {
		b.sessionHandlers[sessionType][sessionState] = HandlerConfig{
			handler:     handler,
			middlewares: middlewares,
		}
	} else {
		log.Fatalf("type: %s. state: %s handler already subscribed", sessionType, sessionState)
	}
}

// Publish handles specific session type and state with middleware chain execution
func (b *TelegramBot) Publish(sessionType BotSessionType, sessionState BotState, c telebot.Context) {
	if handlers, found := b.sessionHandlers[sessionType]; found {
		if handlerConfig, found := handlers[sessionState]; found {
			handler := handlerConfig.handler

			// Apply middlewares in reverse order
			for i := len(handlerConfig.middlewares) - 1; i >= 0; i-- {
				handler = handlerConfig.middlewares[i](handler)
			}

			go handler(c)
		}
	}
}
