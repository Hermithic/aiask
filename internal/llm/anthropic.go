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
	client             anthropic.Client
	model              string
	systemPromptSuffix string
}

// NewAnthropic creates a new Anthropic provider
func NewAnthropic(apiKey, model, systemPromptSuffix string) (*Anthropic, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for Anthropic")
	}

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &Anthropic{
		client:             client,
		model:              model,
		systemPromptSuffix: systemPromptSuffix,
	}, nil
}

// GenerateCommand generates a shell command using Anthropic's Claude API
func (a *Anthropic) GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error) {
	systemPrompt := BuildSmartSystemPrompt(shellInfo, a.systemPromptSuffix, prompt)

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

// ExplainCommand explains what a shell command does
func (a *Anthropic) ExplainCommand(ctx context.Context, command string) (string, error) {
	systemPrompt := BuildExplainPrompt()

	resp, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(a.model),
		MaxTokens: 1000,
		System: []anthropic.TextBlockParam{
			{
				Text: systemPrompt,
				Type: "text",
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(command)),
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

// GenerateCommandStream generates a shell command with streaming output
func (a *Anthropic) GenerateCommandStream(ctx context.Context, prompt string, shellInfo shell.ShellInfo, callback func(chunk string)) (string, error) {
	systemPrompt := BuildSmartSystemPrompt(shellInfo, a.systemPromptSuffix, prompt)

	stream := a.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
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

	var fullContent string
	for stream.Next() {
		event := stream.Current()

		// Check for text delta events
		switch eventVariant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				fullContent += deltaVariant.Text
				if callback != nil {
					callback(deltaVariant.Text)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("streaming API request failed: %w", err)
	}

	return fullContent, nil
}