package llm

import (
	"context"
	"fmt"

	"github.com/Hermithic/aiask/internal/shell"
	"google.golang.org/genai"
)

// Gemini is a provider for Google's Gemini API
type Gemini struct {
	client             *genai.Client
	model              string
	systemPromptSuffix string
}

// NewGemini creates a new Gemini provider
func NewGemini(apiKey, model, systemPromptSuffix string) (*Gemini, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for Gemini")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Gemini{
		client:             client,
		model:              model,
		systemPromptSuffix: systemPromptSuffix,
	}, nil
}

// GenerateCommand generates a shell command using Google's Gemini API
func (g *Gemini) GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error) {
	systemPrompt := BuildSystemPrompt(shellInfo, g.systemPromptSuffix)
	fullPrompt := systemPrompt + "\n\nUser request: " + prompt

	resp, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(fullPrompt), nil)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	// Get the text from the first part
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			return part.Text, nil
		}
	}

	return "", fmt.Errorf("no text response from API")
}

// ExplainCommand explains what a shell command does
func (g *Gemini) ExplainCommand(ctx context.Context, command string) (string, error) {
	systemPrompt := BuildExplainPrompt()
	fullPrompt := systemPrompt + "\n\nCommand to explain: " + command

	resp, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(fullPrompt), nil)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	// Get the text from the first part
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			return part.Text, nil
		}
	}

	return "", fmt.Errorf("no text response from API")
}
