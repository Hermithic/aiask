package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/Hermithic/aiask/internal/config"
	appcontext "github.com/Hermithic/aiask/internal/context"
	"github.com/Hermithic/aiask/internal/shell"
)

// Provider is the interface for LLM providers
type Provider interface {
	// GenerateCommand generates a shell command from a natural language prompt
	GenerateCommand(ctx context.Context, prompt string, shellInfo shell.ShellInfo) (string, error)
	// ExplainCommand explains what a shell command does
	ExplainCommand(ctx context.Context, command string) (string, error)
}

// Closeable is an optional interface for providers that hold resources
type Closeable interface {
	Close() error
}

// CloseProvider closes the provider if it implements Closeable
func CloseProvider(p Provider) error {
	if closer, ok := p.(Closeable); ok {
		return closer.Close()
	}
	return nil
}

// StreamingProvider is an optional interface for providers that support streaming
type StreamingProvider interface {
	Provider
	// GenerateCommandStream generates a shell command with streaming output
	GenerateCommandStream(ctx context.Context, prompt string, shellInfo shell.ShellInfo, callback func(chunk string)) (string, error)
}

// SupportsStreaming checks if a provider supports streaming
func SupportsStreaming(p Provider) bool {
	_, ok := p.(StreamingProvider)
	return ok
}

// NewProvider creates a new LLM provider based on the configuration
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case config.ProviderGrok:
		return NewOpenAICompatible(cfg.APIKey, config.GetProviderURL(cfg.Provider, ""), cfg.Model, cfg.SystemPromptSuffix)
	case config.ProviderOpenAI:
		return NewOpenAICompatible(cfg.APIKey, config.GetProviderURL(cfg.Provider, ""), cfg.Model, cfg.SystemPromptSuffix)
	case config.ProviderOllama:
		return NewOpenAICompatible("", config.GetProviderURL(cfg.Provider, cfg.OllamaURL), cfg.Model, cfg.SystemPromptSuffix)
	case config.ProviderAnthropic:
		return NewAnthropic(cfg.APIKey, cfg.Model, cfg.SystemPromptSuffix)
	case config.ProviderGemini:
		return NewGemini(cfg.APIKey, cfg.Model, cfg.SystemPromptSuffix)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}

// BuildSystemPrompt builds the system prompt for the LLM
func BuildSystemPrompt(shellInfo shell.ShellInfo, suffix string) string {
	prompt := fmt.Sprintf(`You are a shell command assistant. Given a natural language request, return ONLY the shell command(s) needed to accomplish the task.

Rules:
- Return ONLY the command(s), no explanations, no markdown, no code blocks
- If multiple commands are needed, put each on a new line
- Use the appropriate syntax for the current shell
- For dangerous operations, include appropriate safety flags when possible

Current shell: %s
Operating system: %s
Current directory: %s`, shell.GetShellName(shellInfo.Shell), shell.GetOSName(), appcontext.GetCWD())

	// Add git context if in a git repository
	gitCtx := appcontext.GetGitContext()
	if gitCtx.IsRepo {
		prompt += fmt.Sprintf("\nGit branch: %s", gitCtx.Branch)
		if gitCtx.IsDirty {
			prompt += " (has uncommitted changes)"
		}
	}

	if suffix != "" {
		prompt += "\n\nAdditional instructions:\n" + suffix
	}

	return prompt
}

// BuildSystemPromptWithDirContext builds the system prompt with directory listing
func BuildSystemPromptWithDirContext(shellInfo shell.ShellInfo, suffix string) string {
	prompt := BuildSystemPrompt(shellInfo, suffix)
	dirContext := appcontext.GetDirectoryContext()
	if dirContext != "" {
		prompt += "\n\n" + dirContext
	}
	return prompt
}

// BuildSmartSystemPrompt builds the system prompt with context tailored to the user's request
// It includes directory or git context only when relevant to the prompt
func BuildSmartSystemPrompt(shellInfo shell.ShellInfo, suffix string, userPrompt string) string {
	prompt := BuildSystemPrompt(shellInfo, suffix)

	// Add directory context if the prompt is file-related
	if appcontext.IsFileRelatedPrompt(userPrompt) {
		dirContext := appcontext.GetDirectoryContext()
		if dirContext != "" {
			prompt += "\n\n" + dirContext
		}
	}

	// Add extended git context if the prompt is git-related
	if appcontext.IsGitRelatedPrompt(userPrompt) {
		gitStatus := appcontext.GetGitStatus()
		if gitStatus != "" {
			prompt += "\n\n" + gitStatus
		}
		// Add recent commits for context
		recentCommits := appcontext.GetRecentCommits(5)
		if recentCommits != "" {
			prompt += "\nRecent commits:\n" + recentCommits
		}
	}

	return prompt
}

// BuildExplainPrompt builds the system prompt for explaining commands
func BuildExplainPrompt() string {
	return `You are a shell command explainer. Given a shell command, explain what it does in plain English.

Rules:
- Break down each part of the command
- Explain flags and options
- Mention any potential risks or side effects
- Keep explanations clear and concise`
}

// CleanCommand removes markdown code blocks and extra whitespace from the command
func CleanCommand(command string) string {
	// Remove markdown code blocks
	command = strings.TrimPrefix(command, "```bash\n")
	command = strings.TrimPrefix(command, "```powershell\n")
	command = strings.TrimPrefix(command, "```cmd\n")
	command = strings.TrimPrefix(command, "```shell\n")
	command = strings.TrimPrefix(command, "```sh\n")
	command = strings.TrimPrefix(command, "```\n")
	command = strings.TrimSuffix(command, "\n```")
	command = strings.TrimSuffix(command, "```")

	// Trim whitespace
	command = strings.TrimSpace(command)

	return command
}

