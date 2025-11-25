# 🤖 AIask

> **Turn plain English into shell commands instantly!**

[![Release](https://img.shields.io/github/v/release/Hermithic/aiask)](https://github.com/Hermithic/aiask/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8.svg)](https://golang.org/)

AIask is a powerful command-line assistant that understands what you want to do and gives you the exact shell command. No more googling syntax or reading man pages!

```
$ aiask "find all files larger than 100MB"

✨ Suggested command:
   find . -type f -size +100M

What would you like to do?
  [e]xecute  |  [c]opy  |  e[d]it  |  [r]e-prompt  |  [q]uit
> 
```

---

## ✨ Features

### Core Features
- 🗣️ **Natural Language** — Just describe what you want in plain English
- 🐚 **Multi-Shell** — Works with PowerShell, CMD, Bash, Zsh, and Fish
- 🧠 **Multiple AI Providers** — Grok, OpenAI, Anthropic, Google Gemini, or local Ollama
- ⚡ **Interactive** — Execute, copy, edit, or refine commands before running
- 🖥️ **Cross-Platform** — Windows, macOS, and Linux

### New in v2.0 🎉
- 📜 **Command History** — Track and search your command history
- 📝 **Templates** — Save and reuse frequently used prompts
- 🔍 **Explain Mode** — Understand what any command does
- 💬 **Interactive REPL** — Continuous conversation mode
- 🛡️ **Safety Warnings** — Detect and warn about dangerous commands
- ↩️ **Undo Suggestions** — Get undo commands after execution
- 🔧 **Error Recovery** — Get help when commands fail
- 📥 **Stdin Support** — Pipe output for analysis
- 🎨 **Syntax Highlighting** — Colorized command output
- 🌍 **Environment Variables** — Configure via env vars for CI/CD
- 📤 **JSON Output** — Machine-readable output for scripting
- 🐛 **Verbose Mode** — Debug information when needed
- ⏱️ **Configurable Timeout** — Adjust request timeouts
- 🔄 **Auto-Update Check** — Know when updates are available
- 🎯 **Shell Completions** — Tab completion for Bash, Zsh, Fish, PowerShell

---

## 📦 Installation

### 🪟 Windows

**Option 1: winget**
```powershell
winget install Hermithic.aiask
```

**Option 2: Direct download**
```powershell
# Download the latest release
Invoke-WebRequest -Uri "https://github.com/Hermithic/aiask/releases/latest/download/aiask-2.0.0-windows-amd64.zip" -OutFile aiask.zip
Expand-Archive aiask.zip -DestinationPath .
Move-Item aiask-windows-amd64.exe C:\Windows\aiask.exe
```

### 🍎 macOS

```bash
brew tap Hermithic/aiask
brew install aiask
```

### 🐧 Linux

**Option 1: APT (Debian/Ubuntu)**
```bash
# Add the repository
echo "deb [trusted=yes] https://hermithic.github.io/aiask/ stable main" | sudo tee /etc/apt/sources.list.d/aiask.list

# Install
sudo apt update
sudo apt install aiask
```

**Option 2: Download .deb package**
```bash
wget https://github.com/Hermithic/aiask/releases/latest/download/aiask_2.0.0_amd64.deb
sudo dpkg -i aiask_2.0.0_amd64.deb
```

**Option 3: Direct binary**
```bash
wget https://github.com/Hermithic/aiask/releases/latest/download/aiask-2.0.0-linux-amd64.tar.gz
tar -xzf aiask-2.0.0-linux-amd64.tar.gz
sudo mv aiask-linux-amd64 /usr/local/bin/aiask
```

### 🔧 From Source

Requires Go 1.23+:
```bash
git clone https://github.com/Hermithic/aiask.git
cd aiask
go build -o aiask ./cmd/aiask
```

---

## ⚙️ Configuration

### Quick Setup

```bash
aiask config
```

This interactive wizard helps you configure your AI provider.

### Supported Providers

| Provider | Description | API Key Required |
|----------|-------------|------------------|
| 🚀 **Grok** | xAI's Grok (recommended) | Yes |
| 🤖 **OpenAI** | GPT-4o, GPT-4 | Yes |
| 🧠 **Anthropic** | Claude 3.5/4 | Yes |
| ✨ **Gemini** | Google Gemini | Yes |
| 🏠 **Ollama** | Run locally, free! | No |

### 🔑 Getting API Keys

| Provider | Where to get it |
|----------|-----------------|
| Grok | [console.x.ai](https://console.x.ai/) |
| OpenAI | [platform.openai.com](https://platform.openai.com/api-keys) |
| Anthropic | [console.anthropic.com](https://console.anthropic.com/) |
| Gemini | [ai.google.dev](https://ai.google.dev/) |
| Ollama | No key needed! [ollama.ai](https://ollama.ai/) |

### 📁 Configuration File

AIask stores your settings in `~/.aiask/config.yaml`:

```yaml
provider: grok
api_key: "xai-..."
model: "grok-3"
timeout: 60                    # Request timeout in seconds
ollama_url: "http://localhost:11434"
system_prompt_suffix: ""       # Custom instructions for the AI
check_updates: true            # Check for updates on startup
```

### 🌍 Environment Variables

Configure AIask without a config file (great for CI/CD):

```bash
export AIASK_PROVIDER=openai
export AIASK_API_KEY=sk-...
export AIASK_MODEL=gpt-4o
export AIASK_TIMEOUT=120
export AIASK_OLLAMA_URL=http://localhost:11434
export AIASK_SYSTEM_PROMPT_SUFFIX="Prefer one-liners when possible"
```

> Environment variables take precedence over the config file.

### 🏠 Using Ollama (100% Local & Free)

For maximum privacy, run AI completely locally:

```bash
# 1. Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# 2. Pull a model
ollama pull llama3.2

# 3. Configure AIask
aiask config  # Select "ollama"
```

---

## 🚀 Usage

### Basic Usage

```bash
aiask "your request in plain English"
```

### 📚 Examples

**File Operations:**
```bash
aiask "list all .txt files modified in the last 7 days"
aiask "find and delete all empty directories"
aiask "compress the logs folder into a zip"
aiask "count lines of code in all Python files"
```

**System Info:**
```bash
aiask "show disk usage for each partition"
aiask "list all running processes sorted by memory"
aiask "what's my public IP address"
aiask "show system uptime"
```

**Git Operations:**
```bash
aiask "undo my last commit but keep the changes"
aiask "show commits from the last week"
aiask "create a new branch called feature-login"
aiask "squash the last 3 commits"
```

**Networking:**
```bash
aiask "list all open ports"
aiask "download this file from URL"
aiask "check if google.com is reachable"
aiask "show my network interfaces"
```

### ⌨️ Interactive Options

After AIask suggests a command:

| Key | Action |
|-----|--------|
| `e` | ▶️ **Execute** — Run the command now |
| `c` | 📋 **Copy** — Copy to clipboard |
| `d` | ✏️ **Edit** — Modify before running |
| `r` | 🔄 **Re-prompt** — Ask something different |
| `q` | 👋 **Quit** — Exit without action |

---

## 🆕 New Features in v2.0

### 🔍 Explain Mode

Understand what any command does:

```bash
aiask explain "tar -xzvf archive.tar.gz"
aiask explain "git rebase -i HEAD~3"
aiask explain "find . -name '*.log' -mtime +7 -delete"
```

### 📜 Command History

Track and search your command history:

```bash
aiask history              # Show recent history
aiask history -n 20        # Show last 20 entries
aiask history --search git # Search history
aiask history clear        # Clear all history
```

### 📝 Templates

Save and reuse frequently used prompts:

```bash
# Save a template
aiask save git-log "show commits from the last week with stats"
aiask save find-large "find files larger than 100MB" -d "Find large files"

# List templates
aiask templates

# Run a template
aiask run git-log
```

### 💬 Interactive REPL Mode

Continuous conversation mode without restarting:

```bash
aiask interactive
# or
aiask i
```

Commands in REPL:
- `/help` — Show available commands
- `/history` — Show session history
- `/config` — Show current configuration
- `/clear` — Clear the screen
- `/exit` — Exit interactive mode

### 📥 Stdin Support

Pipe output for analysis:

```bash
# Analyze error logs
cat error.log | aiask --stdin "what's wrong here?"

# Get help with failed commands
npm install 2>&1 | aiask --stdin "how do I fix this?"

# Analyze any output
docker logs myapp | aiask --stdin "find any errors"
```

### 🛡️ Safety Features

AIask automatically warns about dangerous commands:

```bash
$ aiask "delete all files in root"

⚠️  CRITICAL Warning
   • Recursive delete of root, all files, or home directory

   Type 'yes' to confirm execution, or any other key to cancel.
```

After execution, get undo suggestions:

```
💡 To undo: git reset HEAD~1
   (Undo the last commit, keeps changes staged)
```

### 📤 JSON Output

Machine-readable output for scripting:

```bash
aiask --json "list files" | jq .command
```

Output:
```json
{
  "command": "ls -la",
  "shell": "bash",
  "os": "linux",
  "prompt": "list files",
  "provider": "grok",
  "model": "grok-3"
}
```

### 🐛 Verbose Mode

Debug information when needed:

```bash
aiask -v "show disk space"
```

Output:
```
[DEBUG] Shell: PowerShell
[DEBUG] OS: Windows
[DEBUG] Provider: grok
[DEBUG] Model: grok-3
[DEBUG] Timeout: 1m0s
[DEBUG] Prompt: show disk space
[DEBUG] Response time: 1.234s
```

### 🎯 Shell Completions

Enable tab completion for your shell:

```bash
# Bash
aiask completion bash > /etc/bash_completion.d/aiask

# Zsh
aiask completion zsh > "${fpath[1]}/_aiask"

# Fish
aiask completion fish > ~/.config/fish/completions/aiask.fish

# PowerShell
aiask completion powershell | Out-String | Invoke-Expression
```

---

## 🤖 Supported Models

| Provider | Default | Other Options |
|----------|---------|---------------|
| Grok | grok-3 | grok-2 |
| OpenAI | gpt-4o | gpt-4o-mini, gpt-4-turbo |
| Anthropic | claude-sonnet-4-20250514 | claude-3-opus |
| Gemini | gemini-2.0-flash | gemini-1.5-pro |
| Ollama | llama3.2 | mistral, codellama, phi |

---

## 🐚 Shell Detection

AIask automatically detects your shell and generates appropriate commands:

| Platform | Shells Detected |
|----------|-----------------|
| Windows | PowerShell, CMD |
| macOS/Linux | Bash, Zsh, Fish |

It also detects:
- Current working directory
- Git repository status (branch, dirty state)

---

## 📋 Command Reference

```
Usage:
  aiask [prompt] [flags]
  aiask [command]

Available Commands:
  config      Configure aiask settings
  explain     Explain what a command does
  history     View command history
  interactive Start interactive REPL mode
  templates   Manage saved prompt templates
  save        Save a new template
  run         Run a saved template
  completion  Generate shell completion scripts
  version     Print the version number
  help        Help about any command

Flags:
  -v, --verbose   Show verbose output including debug information
      --json      Output result as JSON (non-interactive)
      --stdin     Read additional context from stdin
  -s, --stream    Stream the response as it generates
  -h, --help      Help for aiask
```

---

## 🔨 Building from Source

### Prerequisites
- Go 1.23+
- Make (optional)

### Build Commands

```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make release      # Create release archives
make deb          # Build .deb package
make checksums    # Generate SHA256 checksums
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

---

## 🔧 Troubleshooting

### ❌ "Config not found" Error
```bash
aiask config  # Run the setup wizard
```

Or set environment variables:
```bash
export AIASK_PROVIDER=openai
export AIASK_API_KEY=sk-...
```

### ❌ API Errors
1. ✅ Verify your API key is correct
2. ✅ Check you have credits/quota with your provider
3. ✅ Ensure you're using a valid model name
4. ✅ Try increasing the timeout: `export AIASK_TIMEOUT=120`

### ❌ Ollama Connection Issues
```bash
ollama serve    # Make sure Ollama is running
ollama list     # Verify you have models installed
```

### ❌ Command Not Executing
Some commands need elevated privileges:
- **Windows:** Run terminal as Administrator
- **Linux/macOS:** Use `sudo`

### ❌ Slow Responses
Try a faster model:
```yaml
# In ~/.aiask/config.yaml
model: gpt-4o-mini  # Faster than gpt-4o
```

---

## 🔒 Privacy & Security

- 🔐 API keys stored locally in `~/.aiask/config.yaml` with restricted permissions
- 📤 Prompts are sent to your configured AI provider
- 🏠 Use Ollama for 100% local, private inference
- ✋ Commands never execute without your confirmation
- 🛡️ Dangerous commands require explicit "yes" confirmation
- 📜 History stored locally in `~/.aiask/history.yaml`

---

## 🤝 Contributing

Contributions are welcome! Feel free to:
- 🐛 Report bugs
- 💡 Suggest features
- 🔧 Submit pull requests

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

---

## 📝 Changelog

### v2.0.0 (2025-11-26)
- ✨ Added `explain` command to describe what commands do
- ✨ Added command history with search (`aiask history`)
- ✨ Added templates system (`aiask save`, `aiask run`, `aiask templates`)
- ✨ Added interactive REPL mode (`aiask interactive`)
- ✨ Added stdin support for piping input (`--stdin`)
- ✨ Added JSON output mode (`--json`)
- ✨ Added verbose/debug mode (`-v`)
- ✨ Added dangerous command detection and warnings
- ✨ Added undo suggestions after command execution
- ✨ Added error recovery assistance
- ✨ Added shell completion scripts (`aiask completion`)
- ✨ Added environment variable configuration
- ✨ Added configurable timeout
- ✨ Added custom system prompt suffix
- ✨ Added auto-update check on startup
- ✨ Added git context awareness (branch, dirty status)
- ✨ Added directory context in prompts
- ✨ Added syntax highlighting for commands
- 🐛 Fixed various shell detection issues
- 📚 Comprehensive documentation update

### v1.0.0 (2025-01-15)
- 🎉 Initial release
- Basic natural language to command conversion
- Support for multiple AI providers
- Cross-platform support

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/Hermithic">Hermithic</a>
</p>
