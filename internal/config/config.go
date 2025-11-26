package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Hermithic/aiask/internal/fileutil"
	"gopkg.in/yaml.v3"
)

// Environment variable names
const (
	EnvProvider           = "AIASK_PROVIDER"
	EnvAPIKey             = "AIASK_API_KEY"
	EnvModel              = "AIASK_MODEL"
	EnvOllamaURL          = "AIASK_OLLAMA_URL"
	EnvTimeout            = "AIASK_TIMEOUT"
	EnvSystemPromptSuffix = "AIASK_SYSTEM_PROMPT_SUFFIX"
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
	Provider           Provider `yaml:"provider"`
	APIKey             string   `yaml:"api_key,omitempty"`
	Model              string   `yaml:"model"`
	OllamaURL          string   `yaml:"ollama_url,omitempty"`
	Timeout            int      `yaml:"timeout,omitempty"`             // Timeout in seconds (default: 60)
	SystemPromptSuffix string   `yaml:"system_prompt_suffix,omitempty"` // Custom suffix for system prompt
	CheckUpdates       bool     `yaml:"check_updates,omitempty"`       // Whether to check for updates on startup
}

// GetTimeout returns the timeout duration
func (c *Config) GetTimeout() time.Duration {
	if c.Timeout <= 0 {
		return 60 * time.Second
	}
	return time.Duration(c.Timeout) * time.Second
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Provider:     ProviderGrok,
		Model:        "grok-3",
		OllamaURL:    "http://localhost:11434",
		Timeout:      60,
		CheckUpdates: true,
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

// Load loads the configuration from the config file and applies environment variable overrides
func Load() (*Config, error) {
	// Check if we can load entirely from environment variables
	if envProvider := os.Getenv(EnvProvider); envProvider != "" {
		cfg := loadFromEnv()
		if cfg.Provider != "" && (cfg.APIKey != "" || cfg.Provider == ProviderOllama) {
			return cfg, nil
		}
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try loading from env vars only
		cfg := loadFromEnv()
		if cfg.Provider != "" && (cfg.APIKey != "" || cfg.Provider == ProviderOllama) {
			return cfg, nil
		}
		return nil, fmt.Errorf("config not found. Run 'aiask config' to set up, or set AIASK_PROVIDER and AIASK_API_KEY environment variables")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides (env vars take precedence)
	applyEnvOverrides(cfg)

	return cfg, nil
}

// loadFromEnv creates a config entirely from environment variables
func loadFromEnv() *Config {
	cfg := DefaultConfig()

	if provider := os.Getenv(EnvProvider); provider != "" {
		cfg.Provider = Provider(strings.ToLower(provider))
		cfg.Model = GetDefaultModel(cfg.Provider)
	}

	if apiKey := os.Getenv(EnvAPIKey); apiKey != "" {
		cfg.APIKey = apiKey
	}

	if model := os.Getenv(EnvModel); model != "" {
		cfg.Model = model
	}

	if ollamaURL := os.Getenv(EnvOllamaURL); ollamaURL != "" {
		cfg.OllamaURL = ollamaURL
	}

	if timeout := os.Getenv(EnvTimeout); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Timeout = t
		}
	}

	if suffix := os.Getenv(EnvSystemPromptSuffix); suffix != "" {
		cfg.SystemPromptSuffix = suffix
	}

	return cfg
}

// applyEnvOverrides applies environment variable overrides to an existing config
func applyEnvOverrides(cfg *Config) {
	if provider := os.Getenv(EnvProvider); provider != "" {
		cfg.Provider = Provider(strings.ToLower(provider))
	}

	if apiKey := os.Getenv(EnvAPIKey); apiKey != "" {
		cfg.APIKey = apiKey
	}

	if model := os.Getenv(EnvModel); model != "" {
		cfg.Model = model
	}

	if ollamaURL := os.Getenv(EnvOllamaURL); ollamaURL != "" {
		cfg.OllamaURL = ollamaURL
	}

	if timeout := os.Getenv(EnvTimeout); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Timeout = t
		}
	}

	if suffix := os.Getenv(EnvSystemPromptSuffix); suffix != "" {
		cfg.SystemPromptSuffix = suffix
	}
}

// Save saves the configuration to the config file atomically to prevent corruption
func Save(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := fileutil.AtomicWriteFile(configPath, data, 0600); err != nil {
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

