# AIask CLI

**AI-powered command line assistant** that converts natural language into shell commands for PowerShell, CMD, Bash, and Zsh.

```
> aiask "find all files larger than 100MB"

Suggested command:
  find . -type f -size +100M

What would you like to do?
  [e]xecute  |  [c]opy  |  e[d]it  |  [r]e-prompt  |  [q]uit
>
```

## Features

- **Natural Language to Command**: Describe what you want in plain English
- **Multi-Shell Support**: Auto-detects PowerShell, CMD, Bash, Zsh, and Fish
- **Multiple LLM Providers**: Choose from Grok (xAI), OpenAI, Anthropic Claude, Google Gemini, or local Ollama
- **Interactive Workflow**: Execute, copy to clipboard, edit, or refine your request
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Installation

### Windows (winget)

```powershell
winget install Hermithic.aiask
```

### macOS (Homebrew)

```bash
brew tap Hermithic/aiask
brew install aiask
```

### Linux

Download the binary from the [releases page](https://github.com/Hermithic/aiask/releases):

```bash
# Download and install
wget https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-linux-amd64
chmod +x aiask-linux-amd64
sudo mv aiask-linux-amd64 /usr/local/bin/aiask

# Or using the tar.gz
wget https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-amd64.tar.gz
tar -xzf aiask-1.0.0-linux-amd64.tar.gz
chmod +x aiask-linux-amd64
sudo mv aiask-linux-amd64 /usr/local/bin/aiask
```

### From Source

Requires Go 1.21+:

```bash
git clone https://github.com/Hermithic/aiask.git
cd aiask
go build -o aiask ./cmd/aiask
```

## Configuration

Before using AIask, configure your LLM provider:

```bash
aiask config
```

This interactive wizard will guide you through:

1. **Selecting a provider**:
   - `grok` - xAI Grok (recommended)
   - `openai` - OpenAI GPT
   - `anthropic` - Anthropic Claude
   - `gemini` - Google Gemini
   - `ollama` - Local LLM (no API key needed)

2. **Entering your API key** (not required for Ollama)

3. **Selecting a model** (defaults provided)

### Getting API Keys

| Provider | Get API Key |
|----------|-------------|
| xAI Grok | [console.x.ai](https://console.x.ai/) |
| OpenAI | [platform.openai.com](https://platform.openai.com/api-keys) |
| Anthropic | [console.anthropic.com](https://console.anthropic.com/) |
| Google Gemini | [ai.google.dev](https://ai.google.dev/) |
| Ollama | No key needed - [ollama.ai](https://ollama.ai/) |

### Using Ollama (Local LLM)

For privacy-focused users, AIask supports local LLMs via Ollama:

1. Install Ollama: https://ollama.ai/
2. Pull a model: `ollama pull llama3.2`
3. Configure AIask: `aiask config` → select `ollama`

## Usage

### Basic Usage

```bash
aiask "your request in natural language"
```

### Examples

```bash
# File operations
aiask "list all .txt files modified in the last 7 days"
aiask "find and delete all empty directories"
aiask "compress the logs folder into a zip file"

# System info
aiask "show disk usage for each partition"
aiask "list all running processes sorted by memory"
aiask "what's my public IP address"

# Git operations
aiask "undo my last commit but keep the changes"
aiask "show commits from the last week"
aiask "create a new branch and switch to it"

# Network
aiask "list all open ports"
aiask "download a file from this URL"
aiask "check if google.com is reachable"
```

### Interactive Options

After AIask suggests a command, you can:

| Key | Action |
|-----|--------|
| `e` | **Execute** - Run the command immediately |
| `c` | **Copy** - Copy to clipboard |
| `d` | **Edit** - Modify the command before running |
| `r` | **Re-prompt** - Ask a different question |
| `q` | **Quit** - Exit without action |

## Configuration File

AIask stores configuration in `~/.aiask/config.yaml`:

```yaml
provider: grok
api_key: "xai-..."
model: "grok-3"
ollama_url: "http://localhost:11434"
```

### Supported Models

| Provider | Default Model | Other Options |
|----------|---------------|---------------|
| Grok | grok-3 | grok-2 |
| OpenAI | gpt-4o | gpt-4o-mini, gpt-4-turbo |
| Anthropic | claude-sonnet-4-20250514 | claude-3-5-sonnet-20241022, claude-3-opus |
| Gemini | gemini-2.0-flash | gemini-1.5-pro |
| Ollama | llama3.2 | mistral, codellama, phi |

## Shell Detection

AIask automatically detects your current shell and tailors commands accordingly:

- **Windows**: PowerShell, Command Prompt (CMD)
- **macOS/Linux**: Bash, Zsh, Fish

The detection uses environment variables (`PSModulePath`, `SHELL`, `COMSPEC`) to identify the active shell.

## Building from Source

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create release archives
make release

# Build .deb package
make deb
```

### Cross-Compilation

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o aiask.exe ./cmd/aiask

# Linux
GOOS=linux GOARCH=amd64 go build -o aiask ./cmd/aiask

# macOS
GOOS=darwin GOARCH=arm64 go build -o aiask ./cmd/aiask
```

## Troubleshooting

### "Config not found" Error

Run `aiask config` to set up your configuration.

### API Errors

1. Verify your API key is correct
2. Check you have credits/quota with your provider
3. Ensure you're using a valid model name

### Ollama Connection Issues

1. Ensure Ollama is running: `ollama serve`
2. Check the URL in config matches Ollama's address
3. Verify you have a model pulled: `ollama list`

### Command Not Executing

Some commands may require elevated privileges. Try running your terminal as Administrator (Windows) or with `sudo` (Linux/macOS).

## Privacy & Security

- API keys are stored locally in `~/.aiask/config.yaml`
- Your prompts are sent to the configured LLM provider
- For maximum privacy, use Ollama for local inference
- Commands are not executed without your explicit confirmation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
