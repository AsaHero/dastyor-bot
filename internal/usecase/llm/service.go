package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"github.com/AsaHero/dastyor-bot/internal/infrastructure/llm"
	"github.com/sashabaranov/go-openai"
)

type llmService struct {
	contextDeadline time.Duration
	llmAPI          llm.LlmAPI
	chunkSize       int
	rateLimitWait   time.Duration
}

func New(contextDeadline time.Duration, llmAPI llm.LlmAPI) LLM {
	return &llmService{
		contextDeadline: contextDeadline,
		llmAPI:          llmAPI,
		chunkSize:       100,
		rateLimitWait:   time.Microsecond * 200,
	}
}

func (s *llmService) Rewrite(ctx context.Context, text string) (<-chan string, <-chan error) {
	answerChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(answerChan)
		defer close(errChan)

		ctx, cancel := context.WithTimeout(ctx, s.contextDeadline)
		defer cancel()

		outputChan, streamErrChan := s.llmAPI.ChatCompletionStreaming(ctx, "", []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(entity.RewriteTemplate, text),
			},
		})

		var accumulator strings.Builder
		ticker := time.NewTicker(s.rateLimitWait)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case err := <-streamErrChan:
				if err != nil {
					if accumulator.Len() > 0 {
						answerChan <- accumulator.String()
					}

					errChan <- inerr.WithMessage(err, "streaming error:")
				}
				return
			case chunk, ok := <-outputChan:
				if !ok {
					if accumulator.Len() > 0 {
						answerChan <- accumulator.String()
					}
					return
				}

				if chunk != "" {
					accumulator.WriteString(chunk)

					// If accumulated content exceeds chunk size, send it
					if accumulator.Len() >= s.chunkSize {
						<-ticker.C // Wait for rate limi
						answerChan <- accumulator.String()
						accumulator.Reset()
					}
				}
			}
		}
	}()

	return answerChan, errChan
}
