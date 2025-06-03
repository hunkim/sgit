# sgit - Solar LLM-powered Git Wrapper

`sgit` is a command-line tool that wraps git and uses Upstage's Solar LLM to automatically generate meaningful commit messages based on your code changes.

## Features

- 🤖 **AI-powered commit messages**: Automatically generates conventional commit messages using Solar LLM
- 🧠 **Smart file staging**: AI analyzes untracked files to decide what should be added to git
- ⚡ **Easy setup**: Simple configuration with your Upstage API key
- 🔄 **Git compatibility**: Works seamlessly with your existing git workflow
- 🎛️ **Flexible options**: Interactive mode, manual override, and traditional git commit fallback
- 📦 **Single binary**: No dependencies, just download and run
- 🛡️ **Smart filtering**: Automatically skips binary files and large files

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

### Smart Add (New!)

AI-powered file staging that analyzes untracked files:

```bash
# Analyze all untracked files and get AI recommendations
sgit add --all

# Preview what would be added without actually adding
sgit add --all --dry-run

# Add files without AI confirmation (based on file type detection only)
sgit add --all --force

# Traditional behavior: add specific files
sgit add file1.go file2.js
```

The smart add command:
- 📊 **Analyzes file content** using Solar LLM
- 🚫 **Skips binary files** automatically (images, executables, archives, etc.)
- 📏 **Skips large files** (> 1MB)
- 🔍 **Detects sensitive files** (API keys, passwords, etc.)
- 🏗️ **Identifies build artifacts** and temporary files
- ✅ **Recommends source files** for version control

### AI-Powered Commits

1. Stage your changes (use smart add or traditional git add):
   ```bash
   sgit add --all          # Smart AI-powered staging
   # or
   git add .               # Traditional staging
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

### Complete Workflow Examples

```bash
# Complete AI-powered workflow
sgit add --all              # Smart file staging
sgit commit                 # AI-generated commit message

# Preview and review workflow
sgit add --all --dry-run    # See what would be added
sgit add --all              # Confirm and add files
sgit commit -i              # Interactive commit with editable message

# Mixed workflow
sgit add --all              # AI file selection
sgit commit -m "manual message"  # Manual commit message

# Force workflow (no AI, but smart filtering)
sgit add --all --force      # Add all non-binary files
sgit commit --no-ai         # Traditional git commit
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

### Smart Add
1. **File Discovery**: Finds all untracked files using `git ls-files --others --exclude-standard`
2. **Binary Detection**: Skips files based on extension and content analysis
3. **Size Filtering**: Skips files larger than 1MB
4. **AI Analysis**: Sends file content to Solar LLM for decision
5. **User Review**: Shows recommendations with reasons
6. **Execution**: Adds approved files to staging area

### AI Commits
1. **Diff Analysis**: sgit reads your staged changes using `git diff --cached`
2. **AI Processing**: Sends the diff to Solar LLM with a specialized prompt
3. **Message Generation**: Solar LLM generates a conventional commit message
4. **User Review**: You can review, edit, or approve the generated message
5. **Git Commit**: Executes `git commit` with the final message

## File Type Detection

sgit automatically handles these file types:

**✅ Always Added:**
- Source code (.go, .js, .py, .java, .c, .cpp, etc.)
- Configuration files (.json, .yaml, .toml, .xml, etc.)
- Documentation (.md, .txt, .rst, etc.)
- Build files (Makefile, package.json, go.mod, etc.)

**❌ Always Skipped:**
- Binary executables (.exe, .dll, .so, .dylib)
- Images (.jpg, .png, .gif, .ico, etc.)
- Media files (.mp3, .mp4, .avi, etc.)
- Archives (.zip, .tar, .gz, .rar, etc.)
- Large files (> 1MB)

**🤖 AI Analyzed:**
- Unknown file types
- Files that might contain sensitive data
- Generated files that could be build artifacts
- Configuration files that might contain secrets

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
3. Commit your changes (`sgit add --all && sgit commit` 😉)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 [Report bugs](https://github.com/hunkim/sgit/issues)
- 💡 [Request features](https://github.com/hunkim/sgit/issues)
- 📖 [Documentation](https://github.com/hunkim/sgit)

---

**Note**: This tool requires an Upstage API key and makes API calls to generate commit messages and analyze files. Please be mindful of your API usage and associated costs. 