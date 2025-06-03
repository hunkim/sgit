# sgit - Solar LLM-powered Git Wrapper

`sgit` is a command-line tool that wraps git and uses Upstage's Solar LLM to automatically generate meaningful commit messages based on your code changes.

## Features

- 🤖 **AI-powered commit messages**: Automatically generates conventional commit messages using Solar LLM
- ⚡ **Easy setup**: Simple configuration with your Upstage API key
- 🔄 **Git compatibility**: Works seamlessly with your existing git workflow
- 🎛️ **Flexible options**: Interactive mode, manual override, and traditional git commit fallback
- 📦 **Single binary**: No dependencies, just download and run

## Installation

### Option 1: Install with Go (Recommended)

If you have Go installed:

```bash
go install github.com/hunkim/sgit@latest
```

### Option 2: Download Binary

Download the latest binary for your platform from the [releases page](https://github.com/hunkim/sgit/releases).

### Option 3: Build from Source

```bash
git clone https://github.com/hunkim/sgit.git
cd sgit
go build -o sgit
```

## Setup

Before using sgit, you need to configure your Upstage API key:

```bash
sgit config
```

This will prompt you for:
- **Upstage API Key**: Your API key from [Upstage Console](https://console.upstage.ai/)
- **Model Name**: Solar model to use (default: `solar-pro2-preview`)

Configuration is saved to `~/.config/sgit/config.yaml`.

## Usage

### Basic Usage

1. Stage your changes as usual:
   ```bash
   git add .
   ```

2. Commit with AI-generated message:
   ```bash
   sgit commit
   ```

   This will:
   - Analyze your staged changes
   - Generate a commit message using Solar LLM
   - Ask for confirmation before committing

### Advanced Usage

**Interactive Mode** - Review and edit the generated message:
```bash
sgit commit -i
```

**Manual Message** - Skip AI generation:
```bash
sgit commit -m "your commit message"
```

**Skip AI** - Use traditional git commit:
```bash
sgit commit --no-ai
```

### Examples

```bash
# Stage files and commit with AI-generated message
git add src/main.go
sgit commit

# Interactive mode to edit the generated message
sgit commit -i

# Bypass AI and use manual message
sgit commit -m "fix: resolve authentication bug"

# Use traditional git commit interface
sgit commit --no-ai
```

## Configuration

Configuration file location: `~/.config/sgit/config.yaml`

```yaml
upstage_api_key: "up_****************************"
upstage_model_name: "solar-pro2-preview"
```

## API Key Setup

1. Sign up at [Upstage Console](https://console.upstage.ai/)
2. Navigate to API Keys section
3. Create a new API key
4. Run `sgit config` and enter your API key

## How It Works

1. **Diff Analysis**: sgit reads your staged changes using `git diff --cached`
2. **AI Processing**: Sends the diff to Solar LLM with a specialized prompt
3. **Message Generation**: Solar LLM generates a conventional commit message
4. **User Review**: You can review, edit, or approve the generated message
5. **Git Commit**: Executes `git commit` with the final message

## Commit Message Format

sgit generates commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Test additions/modifications
- `chore:` - Maintenance tasks

## Requirements

- Git installed and configured
- Upstage API key
- Internet connection for API calls

## Supported Platforms

- Linux (x86_64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (x86_64)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`sgit commit` 😉)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 [Report bugs](https://github.com/hunkim/sgit/issues)
- 💡 [Request features](https://github.com/hunkim/sgit/issues)
- 📖 [Documentation](https://github.com/hunkim/sgit)

---

**Note**: This tool requires an Upstage API key and makes API calls to generate commit messages. Please be mindful of your API usage and associated costs. 