package llm

import (
	"context"
	"fmt"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/shell"
)

// Provider is the interface for LLM providers
type Provider interface {
	// GenerateCommand generates a shell command from a natural language prompt
	GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error)
}

// NewProvider creates a new LLM provider based on the configuration
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case config.ProviderGrok:
		return NewOpenAICompatible(cfg.APIKey, config.GetProviderURL(cfg.Provider, ""), cfg.Model)
	case config.ProviderOpenAI:
		return NewOpenAICompatible(cfg.APIKey, config.GetProviderURL(cfg.Provider, ""), cfg.Model)
	case config.ProviderOllama:
		return NewOpenAICompatible("", config.GetProviderURL(cfg.Provider, cfg.OllamaURL), cfg.Model)
	case config.ProviderAnthropic:
		return NewAnthropic(cfg.APIKey, cfg.Model)
	case config.ProviderGemini:
		return NewGemini(cfg.APIKey, cfg.Model)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}

// BuildSystemPrompt builds the system prompt for the LLM
func BuildSystemPrompt(shellInfo shell.ShellInfo) string {
	return fmt.Sprintf(`You are a shell command assistant. Given a natural language request, return ONLY the shell command(s) needed to accomplish the task.

Rules:
- Return ONLY the command(s), no explanations, no markdown, no code blocks
- If multiple commands are needed, put each on a new line
- Use the appropriate syntax for the current shell
- For dangerous operations, include appropriate safety flags when possible

Current shell: %s
Operating system: %s`, shell.GetShellName(shellInfo.Shell), shell.GetOSName())
}

