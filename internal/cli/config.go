package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/ui"
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

func runConfig(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════╗%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s║        AIask Configuration Setup         ║%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════╝%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()

	// Load existing config if available
	existingCfg, _ := config.Load()
	if existingCfg == nil {
		existingCfg = config.DefaultConfig()
	}

	// Provider selection
	fmt.Println("Select your LLM provider:")
	fmt.Println()
	providers := []struct {
		name string
		desc string
	}{
		{"grok", "xAI Grok (recommended)"},
		{"openai", "OpenAI GPT"},
		{"anthropic", "Anthropic Claude"},
		{"gemini", "Google Gemini"},
		{"ollama", "Ollama (local, no API key needed)"},
	}

	for i, p := range providers {
		marker := "  "
		if config.Provider(p.name) == existingCfg.Provider {
			marker = "→ "
		}
		fmt.Printf("%s%d. %s - %s%s\n", marker, i+1, p.name, p.desc, ui.ColorReset)
	}

	fmt.Println()
	fmt.Printf("Enter choice [1-5] (current: %s): ", existingCfg.Provider)

	providerInput, _ := reader.ReadString('\n')
	providerInput = strings.TrimSpace(providerInput)

	var selectedProvider config.Provider
	if providerInput == "" {
		selectedProvider = existingCfg.Provider
	} else {
		idx, err := strconv.Atoi(providerInput)
		if err != nil || idx < 1 || idx > len(providers) {
			fmt.Println("Invalid selection. Using default (grok).")
			selectedProvider = config.ProviderGrok
		} else {
			selectedProvider = config.Provider(providers[idx-1].name)
		}
	}

	cfg := &config.Config{
		Provider:  selectedProvider,
		Model:     config.GetDefaultModel(selectedProvider),
		OllamaURL: existingCfg.OllamaURL,
	}

	// API Key (not needed for Ollama)
	if selectedProvider != config.ProviderOllama {
		fmt.Println()
		fmt.Printf("Enter your %s API key: ", selectedProvider)

		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		if apiKey == "" && existingCfg.APIKey != "" {
			apiKey = existingCfg.APIKey
			fmt.Println("Using existing API key.")
		}

		if apiKey == "" {
			fmt.Println("Warning: No API key provided. You'll need to set it later.")
		}

		cfg.APIKey = apiKey
	}

	// Model selection
	fmt.Println()
	fmt.Printf("Enter model name (default: %s): ", cfg.Model)

	modelInput, _ := reader.ReadString('\n')
	modelInput = strings.TrimSpace(modelInput)

	if modelInput != "" {
		cfg.Model = modelInput
	}

	// Ollama URL (only for Ollama)
	if selectedProvider == config.ProviderOllama {
		fmt.Println()
		fmt.Printf("Enter Ollama URL (default: %s): ", config.DefaultConfig().OllamaURL)

		ollamaURL, _ := reader.ReadString('\n')
		ollamaURL = strings.TrimSpace(ollamaURL)

		if ollamaURL != "" {
			cfg.OllamaURL = ollamaURL
		} else {
			cfg.OllamaURL = config.DefaultConfig().OllamaURL
		}
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		fmt.Printf("\n%sError saving configuration: %s%s\n", ui.ColorYellow, err, ui.ColorReset)
		os.Exit(1)
	}

	configPath, _ := config.GetConfigPath()

	fmt.Println()
	fmt.Printf("%s✓ Configuration saved successfully!%s\n", ui.ColorGreen, ui.ColorReset)
	fmt.Printf("  Location: %s\n", configPath)
	fmt.Println()
	fmt.Println("You can now use aiask:")
	fmt.Printf("  %saiask \"list all files in current directory\"%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()
}
