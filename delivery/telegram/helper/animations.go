package helper

import (
	"time"

	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"gopkg.in/telebot.v3"
)

func Animate(msg *telebot.Message, bot *telebot.Bot, frames []string, done chan struct{}, failed chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	idx := 1
	for {
		select {
		case <-ticker.C:
			if _, err := bot.Edit(msg, frames[idx%len(frames)]); err != nil {
				inerr.Err(err)
				return
			}
			idx++
		case <-done:
			return
		case <-failed:
			errorCaption := `Nosozlik yuz berdi! ☹️`
			if _, err := bot.Edit(msg, errorCaption, telebot.ModeHTML); err != nil {
				inerr.Err(err)
				return
			}
			return
		}
	}

}
