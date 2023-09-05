package llm

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type openAI struct {
	client *openai.Client
}

func NewOpenAI(apiKey string) *openAI {
	client := openai.NewClient(apiKey)
	return &openAI{
		client: client,
	}
}

func (llm *openAI) GenerateText(ctx context.Context, prompt string, corpus string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: corpus,
			},
		},
		Temperature: 0.7,
		MaxTokens:   100,
		N:           1,
	}
	resp, err := llm.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return string(resp.Choices[0].Message.Content), nil
}
