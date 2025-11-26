package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/Hermithic/aiask/internal/shell"
	"google.golang.org/genai"
)

// Gemini is a provider for Google's Gemini API
type Gemini struct {
	client             *genai.Client
	model              string
	systemPromptSuffix string
}

// Close releases resources associated with the Gemini client
func (g *Gemini) Close() error {
	// Note: genai.Client doesn't require explicit closing
	return nil
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
	systemPrompt := BuildSmartSystemPrompt(shellInfo, g.systemPromptSuffix, prompt)
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

// GenerateCommandStream generates a shell command with streaming output
func (g *Gemini) GenerateCommandStream(ctx context.Context, prompt string, shellInfo shell.ShellInfo, callback func(chunk string)) (string, error) {
	systemPrompt := BuildSmartSystemPrompt(shellInfo, g.systemPromptSuffix, prompt)
	fullPrompt := systemPrompt + "\n\nUser request: " + prompt

	stream := g.client.Models.GenerateContentStream(ctx, g.model, genai.Text(fullPrompt), nil)

	var fullContent string
	for chunk, err := range stream {
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("streaming API request failed: %w", err)
		}

		// Extract text from the chunk
		text, textErr := chunk.Text()
		if textErr != nil {
			continue // Skip chunks without text
		}
		if text != "" {
			fullContent += text
			if callback != nil {
				callback(text)
			}
		}
	}

	return fullContent, nil
}
