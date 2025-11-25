package llm

import (
	"context"
	"fmt"

	"github.com/Hermithic/aiask/internal/shell"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAICompatible is a provider for OpenAI-compatible APIs (OpenAI, Grok, Ollama)
type OpenAICompatible struct {
	client openai.Client
	model  string
}

// NewOpenAICompatible creates a new OpenAI-compatible provider
func NewOpenAICompatible(apiKey, baseURL, model string) (*OpenAICompatible, error) {
	opts := []option.RequestOption{}

	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}

	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(opts...)

	return &OpenAICompatible{
		client: client,
		model:  model,
	}, nil
}

// GenerateCommand generates a shell command using an OpenAI-compatible API
func (o *OpenAICompatible) GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error) {
	systemPrompt := BuildSystemPrompt(shellInfo)

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(o.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}
