package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Provider represents the LLM provider type
type Provider string

const (
	ProviderGrok      Provider = "grok"
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderGemini    Provider = "gemini"
	ProviderOllama    Provider = "ollama"
)

// Config represents the application configuration
type Config struct {
	Provider  Provider `yaml:"provider"`
	APIKey    string   `yaml:"api_key,omitempty"`
	Model     string   `yaml:"model"`
	OllamaURL string   `yaml:"ollama_url,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Provider:  ProviderGrok,
		Model:     "grok-3",
		OllamaURL: "http://localhost:11434",
	}
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".aiask"), nil
}

// GetConfigPath returns the configuration file path
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config not found. Run 'aiask config' to set up")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to the config file
func Save(cfg *Config) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if the configuration file exists
func Exists() bool {
	configPath, err := GetConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(configPath)
	return err == nil
}

// GetDefaultModel returns the default model for a given provider
func GetDefaultModel(provider Provider) string {
	switch provider {
	case ProviderGrok:
		return "grok-3"
	case ProviderOpenAI:
		return "gpt-4o"
	case ProviderAnthropic:
		return "claude-sonnet-4-20250514"
	case ProviderGemini:
		return "gemini-2.0-flash"
	case ProviderOllama:
		return "llama3.2"
	default:
		return ""
	}
}

// GetProviderURL returns the API URL for a given provider
func GetProviderURL(provider Provider, ollamaURL string) string {
	switch provider {
	case ProviderGrok:
		return "https://api.x.ai/v1"
	case ProviderOpenAI:
		return "https://api.openai.com/v1"
	case ProviderOllama:
		if ollamaURL != "" {
			return ollamaURL + "/v1"
		}
		return "http://localhost:11434/v1"
	default:
		return ""
	}
}

// ValidProviders returns a list of valid provider names
func ValidProviders() []string {
	return []string{
		string(ProviderGrok),
		string(ProviderOpenAI),
		string(ProviderAnthropic),
		string(ProviderGemini),
		string(ProviderOllama),
	}
}

