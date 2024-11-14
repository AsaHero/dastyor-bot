package llm

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AsaHero/dastyor-bot/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type apiClient struct {
	client *openai.Client
}

func New(cfg *config.Config) (LlmAPI, error) {
	client := openai.NewClient(cfg.LLM.SecretKey)

	return &apiClient{
		client: client,
	}, nil
}

func (c *apiClient) ChatCompletionStreaming(ctx context.Context, system string, messages []openai.ChatCompletionMessage) (<-chan string, <-chan error) {
	output := make(chan string, 1)
	errchan := make(chan error, 1)
	chatCompletionMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: system,
		},
	}

	chatCompletionMessages = append(chatCompletionMessages, messages...)

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Stream:   true,
		Messages: chatCompletionMessages,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		close(output)
		errchan <- fmt.Errorf("(chatgptAPI) ChatCompletionStream error: %v", err)
		close(errchan)
		return output, errchan
	}

	go func() {
		defer close(output)
		defer close(errchan)
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				response, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}

					errchan <- fmt.Errorf("(chatgptAPI), Error while parsing response chunks: %v", err)
					return
				}

				text := response.Choices[0].Delta.Content
				output <- text
			}
		}
	}()

	return output, nil
}
