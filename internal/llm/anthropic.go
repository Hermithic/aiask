package llm

import (
	"context"
	"fmt"

	"github.com/Hermithic/aiask/internal/shell"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Anthropic is a provider for Anthropic's Claude API
type Anthropic struct {
	client anthropic.Client
	model  string
}

// NewAnthropic creates a new Anthropic provider
func NewAnthropic(apiKey, model string) (*Anthropic, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for Anthropic")
	}

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &Anthropic{
		client: client,
		model:  model,
	}, nil
}

// GenerateCommand generates a shell command using Anthropic's Claude API
func (a *Anthropic) GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error) {
	systemPrompt := BuildSystemPrompt(shellInfo)

	resp, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(a.model),
		MaxTokens: 500,
		System: []anthropic.TextBlockParam{
			{
				Text: systemPrompt,
				Type: "text",
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	// Extract text from response
	for _, block := range resp.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("no text response from API")
}
