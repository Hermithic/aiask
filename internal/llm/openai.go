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
	client             openai.Client
	model              string
	systemPromptSuffix string
}

// NewOpenAICompatible creates a new OpenAI-compatible provider
func NewOpenAICompatible(apiKey, baseURL, model, systemPromptSuffix string) (*OpenAICompatible, error) {
	opts := []option.RequestOption{}

	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}

	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(opts...)

	return &OpenAICompatible{
		client:             client,
		model:              model,
		systemPromptSuffix: systemPromptSuffix,
	}, nil
}

// GenerateCommand generates a shell command using an OpenAI-compatible API
func (o *OpenAICompatible) GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error) {
	systemPrompt := BuildSmartSystemPrompt(shellInfo, o.systemPromptSuffix, prompt)

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(o.model),
		MaxTokens: openai.Int(500),
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

// ExplainCommand explains what a shell command does
func (o *OpenAICompatible) ExplainCommand(ctx context.Context, command string) (string, error) {
	systemPrompt := BuildExplainPrompt()

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(o.model),
		MaxTokens: openai.Int(1000),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(command),
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

// GenerateCommandStream generates a shell command with streaming output
func (o *OpenAICompatible) GenerateCommandStream(ctx context.Context, prompt string, shellInfo shell.ShellInfo, callback func(chunk string)) (string, error) {
	systemPrompt := BuildSmartSystemPrompt(shellInfo, o.systemPromptSuffix, prompt)

	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(o.model),
		MaxTokens: openai.Int(500),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		},
	})

	var fullContent string
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			content := chunk.Choices[0].Delta.Content
			fullContent += content
			if callback != nil {
				callback(content)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("streaming API request failed: %w", err)
	}

	return fullContent, nil
}
