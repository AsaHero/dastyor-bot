package llm

import "context"

type LLM interface {
	Rewrite(ctx context.Context, text string) (<-chan string, <-chan error)
}