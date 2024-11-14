package handler

import (
	"context"

	"github.com/AsaHero/dastyor-bot/delivery/telegram/helper"
	"github.com/AsaHero/dastyor-bot/delivery/telegram/models"
	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"gopkg.in/telebot.v3"
)

func (h *BotHandler) RewriteInit(c telebot.Context) error {
	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		keyboard.Row(*models.ButtonHome),
	)

	return c.Send(models.MessageRewriteMenu, keyboard)
}

func (h *BotHandler) RewriteHandle(c telebot.Context) error {
	ctx := context.Background()
	text := c.Text()

	if text == "" {
		h.SessionsAssigner(c, entity.REWRITE, entity.REWRITE_waiting)
		return h.RewriteInit(c)
	}

	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		keyboard.Row(*models.ButtonHome),
	)

	thinkingMsg, err := c.Bot().Send(c.Chat(), models.AnimationThinking[0])
	if err != nil {
		return inerr.Err(err)
	}

	animationDone := make(chan struct{})
	go helper.Animate(thinkingMsg, c.Bot(), models.AnimationThinking, animationDone, animationDone)

	outputChan, errChan := h.llmService.Rewrite(ctx, text)
	var result string

	// Handle first chunk to stop animation before starting updates
	select {
	case chunk := <-outputChan:
		close(animationDone)
		result = chunk
		_, err := c.Bot().Edit(thinkingMsg, result)
		if err != nil {
			return inerr.Err(err)
		}
	case err := <-errChan:
		inerr.WithMessage(err, "error while getting llm api response")
		return c.Send("Nosozlik yuz berdi! ☹️", result, keyboard)

	}

	// Process remaining chunks
	for {
		select {
		case chunk, ok := <-outputChan:
			if !ok {
				_, err := c.Bot().Edit(thinkingMsg, result, keyboard)
				if err != nil {
					return inerr.Err(err)
				}

				return nil
			}

			result += chunk
			_, err := c.Bot().Edit(thinkingMsg, result)
			if err != nil {
				return inerr.Err(err)
			}

		case err, ok := <-errChan:
			if !ok {
				_, err := c.Bot().Edit(thinkingMsg, result, keyboard)
				if err != nil {
					return inerr.Err(err)
				}

				return nil
			}
			inerr.WithMessage(err, "error while getting llm api response")
			return c.Send("Nosozlik yuz berdi! ☹️", result, keyboard)
		}
	}
}
