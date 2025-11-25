# 🤖 AIask

> **Turn plain English into shell commands instantly!**

[![Release](https://img.shields.io/github/v/release/Hermithic/aiask)](https://github.com/Hermithic/aiask/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

AIask is a command-line assistant that understands what you want to do and gives you the exact shell command. No more googling syntax or reading man pages!

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

- 🗣️ **Natural Language** — Just describe what you want in plain English
- 🐚 **Multi-Shell** — Works with PowerShell, CMD, Bash, Zsh, and Fish
- 🧠 **Multiple AI Providers** — Grok, OpenAI, Anthropic, Google Gemini, or local Ollama
- ⚡ **Interactive** — Execute, copy, edit, or refine commands before running
- 🖥️ **Cross-Platform** — Windows, macOS, and Linux

---

## 📦 Installation

### 🪟 Windows

**Option 1: winget** *(coming soon)*
```powershell
winget install Hermithic.aiask
```

**Option 2: Direct download**
```powershell
# Download the latest release
Invoke-WebRequest -Uri "https://github.com/Hermithic/aiask/releases/latest/download/aiask-1.0.0-windows-amd64.zip" -OutFile aiask.zip
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
wget https://github.com/Hermithic/aiask/releases/latest/download/aiask_1.0.0_amd64.deb
sudo dpkg -i aiask_1.0.0_amd64.deb
```

**Option 3: Direct binary**
```bash
wget https://github.com/Hermithic/aiask/releases/latest/download/aiask-1.0.0-linux-amd64.tar.gz
tar -xzf aiask-1.0.0-linux-amd64.tar.gz
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

Before using AIask, set up your AI provider:

```bash
aiask config
```

This interactive wizard helps you:

1. **Choose a provider:**
   | Provider | Description | API Key Required |
   |----------|-------------|------------------|
   | 🚀 Grok | xAI's Grok (recommended) | Yes |
   | 🤖 OpenAI | GPT-4o, GPT-4 | Yes |
   | 🧠 Anthropic | Claude 3.5/4 | Yes |
   | ✨ Gemini | Google Gemini | Yes |
   | 🏠 Ollama | Run locally, free! | No |

2. **Enter your API key** (not needed for Ollama)

3. **Select a model** (defaults provided)

### 🔑 Getting API Keys

| Provider | Where to get it |
|----------|-----------------|
| Grok | [console.x.ai](https://console.x.ai/) |
| OpenAI | [platform.openai.com](https://platform.openai.com/api-keys) |
| Anthropic | [console.anthropic.com](https://console.anthropic.com/) |
| Gemini | [ai.google.dev](https://ai.google.dev/) |
| Ollama | No key needed! [ollama.ai](https://ollama.ai/) |

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

## 📁 Configuration File

AIask stores your settings in `~/.aiask/config.yaml`:

```yaml
provider: grok
api_key: "xai-..."
model: "grok-3"
ollama_url: "http://localhost:11434"  # Only for Ollama
```

### 🤖 Supported Models

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

### ❌ API Errors
1. ✅ Verify your API key is correct
2. ✅ Check you have credits/quota with your provider
3. ✅ Ensure you're using a valid model name

### ❌ Ollama Connection Issues
```bash
ollama serve    # Make sure Ollama is running
ollama list     # Verify you have models installed
```

### ❌ Command Not Executing
Some commands need elevated privileges:
- **Windows:** Run terminal as Administrator
- **Linux/macOS:** Use `sudo`

---

## 🔒 Privacy & Security

- 🔐 API keys stored locally in `~/.aiask/config.yaml`
- 📤 Prompts are sent to your configured AI provider
- 🏠 Use Ollama for 100% local, private inference
- ✋ Commands never execute without your confirmation

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

<p align="center">
  Made with ❤️ by <a href="https://github.com/Hermithic">Hermithic</a>
</p>
