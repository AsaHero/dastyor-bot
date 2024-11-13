package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type LLM interface {
	ChatCompletionStreaming(ctx context.Context, system string, messages []openai.ChatCompletionMessage) (<-chan string, <-chan error)
}
