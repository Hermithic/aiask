package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure aiask settings",
	Long: `Configure your LLM provider and API key for aiask.

Supported providers:
  - grok      : xAI Grok (default) - https://x.ai/
  - openai    : OpenAI GPT - https://platform.openai.com/
  - anthropic : Anthropic Claude - https://console.anthropic.com/
  - gemini    : Google Gemini - https://ai.google.dev/
  - ollama    : Ollama (local) - https://ollama.ai/`,
	Run: runConfig,
}

// providerOption represents a provider choice in the menu
type providerOption struct {
	Name        string
	Description string
	URL         string
	NeedsAPIKey bool
	Icon        string
}

func runConfig(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(ui.Header("AIask Configuration", 44))
	fmt.Println()

	// Load existing config if available
	existingCfg, _ := config.Load()
	if existingCfg == nil {
		existingCfg = config.DefaultConfig()
	}

	// Provider selection with interactive menu
	providers := []providerOption{
		{Name: "grok", Description: "xAI Grok (recommended)", URL: "https://x.ai/", NeedsAPIKey: true, Icon: "⚡"},
		{Name: "openai", Description: "OpenAI GPT", URL: "https://platform.openai.com/", NeedsAPIKey: true, Icon: "🤖"},
		{Name: "anthropic", Description: "Anthropic Claude", URL: "https://console.anthropic.com/", NeedsAPIKey: true, Icon: "🧠"},
		{Name: "gemini", Description: "Google Gemini", URL: "https://ai.google.dev/", NeedsAPIKey: true, Icon: "✨"},
		{Name: "ollama", Description: "Ollama (local)", URL: "https://ollama.ai/", NeedsAPIKey: false, Icon: "🏠"},
	}

	// Find current provider index
	currentIdx := 0
	for i, p := range providers {
		if config.Provider(p.Name) == existingCfg.Provider {
			currentIdx = i
			break
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   fmt.Sprintf("%s%s {{ .Icon }} {{ .Name | cyan | bold }}%s - {{ .Description }}", ui.ColorCyan, ui.IconArrow, ui.ColorReset),
		Inactive: "  {{ .Icon }} {{ .Name }} - {{ .Description | faint }}",
		Selected: fmt.Sprintf("%s%s {{ .Icon }} {{ .Name }}%s", ui.ColorGreen, ui.IconCheck, ui.ColorReset),
		Details: `
{{ "────────────────────────────────────────────" | faint }}
{{ "Provider:" | faint }}  {{ .Name }}
{{ "API Key:" | faint }}   {{ if .NeedsAPIKey }}Required{{ else }}Not needed (local){{ end }}
{{ "Website:" | faint }}   {{ .URL | faint }}`,
	}

	providerPrompt := promptui.Select{
		Label:     fmt.Sprintf("%sSelect your LLM provider%s", ui.ColorBold, ui.ColorReset),
		Items:     providers,
		Templates: templates,
		Size:      5,
		CursorPos: currentIdx,
	}

	idx, _, err := providerPrompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("\nConfiguration cancelled.")
			return
		}
		ui.ShowError(fmt.Errorf("provider selection failed: %w", err))
		return
	}

	selectedProvider := config.Provider(providers[idx].Name)

	cfg := &config.Config{
		Provider:           selectedProvider,
		Model:              config.GetDefaultModel(selectedProvider),
		OllamaURL:          existingCfg.OllamaURL,
		Timeout:            existingCfg.Timeout,
		SystemPromptSuffix: existingCfg.SystemPromptSuffix,
		CheckUpdates:       existingCfg.CheckUpdates,
	}

	// API Key (not needed for Ollama)
	if providers[idx].NeedsAPIKey {
		fmt.Println()

		hasExistingKey := existingCfg.APIKey != "" && existingCfg.Provider == selectedProvider
		defaultText := ""
		if hasExistingKey {
			maskedKey := maskAPIKey(existingCfg.APIKey)
			defaultText = fmt.Sprintf(" (current: %s)", maskedKey)
		}

		apiKeyPrompt := promptui.Prompt{
			Label: fmt.Sprintf("Enter %s API key%s", selectedProvider, defaultText),
			Mask:  '*',
			Templates: &promptui.PromptTemplates{
				Prompt:  fmt.Sprintf("%s{{ . }}:%s ", ui.ColorCyan, ui.ColorReset),
				Valid:   fmt.Sprintf("%s{{ . }}:%s ", ui.ColorGreen, ui.ColorReset),
				Invalid: fmt.Sprintf("%s{{ . }}:%s ", ui.ColorRed, ui.ColorReset),
				Success: fmt.Sprintf("%s%s {{ . }}:%s ", ui.ColorGreen, ui.IconCheck, ui.ColorReset),
			},
		}

		apiKey, err := apiKeyPrompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Println("\nConfiguration cancelled.")
				return
			}
		}

		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" && hasExistingKey {
			cfg.APIKey = existingCfg.APIKey
			fmt.Printf("%sUsing existing API key.%s\n", ui.ColorDim, ui.ColorReset)
		} else if apiKey == "" {
			fmt.Println(ui.WarningMessage("No API key provided. You'll need to set it later."))
		} else {
			cfg.APIKey = apiKey
		}
	}

	// Model selection
	fmt.Println()

	modelPrompt := promptui.Prompt{
		Label:   "Model name",
		Default: cfg.Model,
		Templates: &promptui.PromptTemplates{
			Prompt:  fmt.Sprintf("%s{{ . }}:%s ", ui.ColorCyan, ui.ColorReset),
			Valid:   fmt.Sprintf("%s{{ . }}:%s ", ui.ColorGreen, ui.ColorReset),
			Invalid: fmt.Sprintf("%s{{ . }}:%s ", ui.ColorRed, ui.ColorReset),
			Success: fmt.Sprintf("%s%s {{ . }}:%s ", ui.ColorGreen, ui.IconCheck, ui.ColorReset),
		},
	}

	modelInput, err := modelPrompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("\nConfiguration cancelled.")
			return
		}
	}

	if strings.TrimSpace(modelInput) != "" {
		cfg.Model = modelInput
	}

	// Ollama URL (only for Ollama)
	if selectedProvider == config.ProviderOllama {
		fmt.Println()

		defaultURL := config.DefaultConfig().OllamaURL
		if existingCfg.OllamaURL != "" {
			defaultURL = existingCfg.OllamaURL
		}

		ollamaPrompt := promptui.Prompt{
			Label:   "Ollama server URL",
			Default: defaultURL,
			Templates: &promptui.PromptTemplates{
				Prompt:  fmt.Sprintf("%s{{ . }}:%s ", ui.ColorCyan, ui.ColorReset),
				Valid:   fmt.Sprintf("%s{{ . }}:%s ", ui.ColorGreen, ui.ColorReset),
				Invalid: fmt.Sprintf("%s{{ . }}:%s ", ui.ColorRed, ui.ColorReset),
				Success: fmt.Sprintf("%s%s {{ . }}:%s ", ui.ColorGreen, ui.IconCheck, ui.ColorReset),
			},
		}

		ollamaURL, err := ollamaPrompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Println("\nConfiguration cancelled.")
				return
			}
		}

		if strings.TrimSpace(ollamaURL) != "" {
			cfg.OllamaURL = ollamaURL
		} else {
			cfg.OllamaURL = defaultURL
		}
	}

	// Confirmation
	fmt.Println()
	fmt.Println(ui.Divider(44))
	fmt.Printf("%sConfiguration Summary:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Printf("  Provider: %s%s%s\n", ui.ColorCyan, cfg.Provider, ui.ColorReset)
	fmt.Printf("  Model:    %s%s%s\n", ui.ColorCyan, cfg.Model, ui.ColorReset)
	if cfg.APIKey != "" {
		fmt.Printf("  API Key:  %s%s%s\n", ui.ColorCyan, maskAPIKey(cfg.APIKey), ui.ColorReset)
	}
	if selectedProvider == config.ProviderOllama {
		fmt.Printf("  URL:      %s%s%s\n", ui.ColorCyan, cfg.OllamaURL, ui.ColorReset)
	}
	fmt.Println(ui.Divider(44))
	fmt.Println()

	// Confirm save
	confirmPrompt := promptui.Prompt{
		Label:     "Save configuration",
		IsConfirm: true,
	}

	_, err = confirmPrompt.Run()
	if err != nil {
		fmt.Println("\nConfiguration not saved.")
		return
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		fmt.Println()
		ui.ShowError(fmt.Errorf("error saving configuration: %w", err))
		os.Exit(1)
	}

	configPath, _ := config.GetConfigPath()

	fmt.Println()
	fmt.Println(ui.SuccessMessage("Configuration saved successfully!"))
	fmt.Printf("  %sLocation:%s %s\n", ui.ColorDim, ui.ColorReset, configPath)
	fmt.Println()
	fmt.Printf("%sYou can now use aiask:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Printf("  %saiask \"list all files in current directory\"%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()
}

// maskAPIKey masks an API key for display
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}
