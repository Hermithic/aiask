# Changelog

All notable changes to AIask will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2025-11-26

### Added
- **Explain Mode**: New `aiask explain` command to describe what any shell command does
- **Command History**: Track and search command history with `aiask history`
  - View recent history: `aiask history`
  - Search history: `aiask history --search git`
  - Clear history: `aiask history clear`
- **Templates System**: Save and reuse frequently used prompts
  - Save template: `aiask save <name> "<prompt>"`
  - Run template: `aiask run <name>`
  - List templates: `aiask templates`
- **Interactive REPL Mode**: Continuous conversation mode with `aiask interactive`
  - Built-in commands: `/help`, `/history`, `/config`, `/clear`, `/exit`
- **Stdin Support**: Pipe output for analysis with `--stdin` flag
  - Example: `cat error.log | aiask --stdin "what's wrong?"`
- **JSON Output Mode**: Machine-readable output with `--json` flag
- **Verbose/Debug Mode**: Debug information with `-v` or `--verbose` flag
- **Dangerous Command Detection**: Automatic warnings for destructive commands
  - Requires explicit "yes" confirmation for dangerous operations
  - Detects: `rm -rf`, `DROP TABLE`, `format`, and more
- **Undo Suggestions**: Shows how to undo commands after execution
- **Error Recovery**: Offers help when commands fail
- **Shell Completions**: Tab completion for Bash, Zsh, Fish, and PowerShell
  - Generate with: `aiask completion <shell>`
- **Environment Variables**: Configure without a config file
  - `AIASK_PROVIDER`, `AIASK_API_KEY`, `AIASK_MODEL`
  - `AIASK_TIMEOUT`, `AIASK_OLLAMA_URL`
  - `AIASK_SYSTEM_PROMPT_SUFFIX`
- **Configurable Timeout**: Set request timeout in config or via env var
- **Custom System Prompts**: Add custom instructions via `system_prompt_suffix`
- **Auto-Update Check**: Notification when updates are available
- **Git Context Awareness**: Detects git branch and dirty status
- **Directory Context**: Includes current working directory in prompts
- **Syntax Highlighting**: Colorized command output
- **Streaming Response**: Stream AI responses with `--stream` flag

### Changed
- Improved shell detection for WSL environments
- Better error messages with suggestions
- Enhanced command cleaning to handle more markdown formats

### Fixed
- Fixed flag conflicts in subcommands
- Fixed history recording for all command actions

## [1.0.0] - 2025-01-15

### Added
- Initial release
- Natural language to shell command conversion
- Support for multiple AI providers:
  - Grok (xAI)
  - OpenAI (GPT-4o, GPT-4)
  - Anthropic (Claude)
  - Google Gemini
  - Ollama (local)
- Interactive command execution options:
  - Execute, Copy, Edit, Re-prompt, Quit
- Auto-detection of current shell (PowerShell, CMD, Bash, Zsh, Fish)
- Cross-platform support (Windows, macOS, Linux)
- Configuration wizard (`aiask config`)
- Multiple installation methods:
  - winget (Windows)
  - Homebrew (macOS)
  - APT/deb (Linux)
  - Direct binary download

[2.0.0]: https://github.com/Hermithic/aiask/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/Hermithic/aiask/releases/tag/v1.0.0

