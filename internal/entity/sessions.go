package entity

import (
	"time"

	telegram_bot "github.com/AsaHero/dastyor-bot/pkg/telegram-bot"
)

const (
	// Rewrite session
	REWRITE           telegram_bot.BotSessionType = "REWRITE"
	REWRITE_waiting   telegram_bot.BotState       = "REWRITE_waiting"
	REWRITE_answering telegram_bot.BotState       = "REWRITE_answering"
)

type Sessions struct {
	ID        int64
	UserID    int64
	Type      telegram_bot.BotSessionType
	State     telegram_bot.BotState
	CreatedAt time.Time
	UpdatedAt time.Time
}
