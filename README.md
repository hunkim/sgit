# sgit - Solar LLM-powered Git Wrapper

`sgit` is a command-line tool that enhances git with AI capabilities while maintaining full compatibility with all existing git commands and workflows. It uses Upstage's Solar LLM to provide intelligent commit messages, smart file staging, diff summaries, log analysis, and merge assistance.

## Key Principles

🔄 **Drop-in Replacement**: sgit works exactly like git - you don't need to learn new commands or change your workflow  
🤖 **Optional AI Enhancement**: AI features are opt-in and never interfere with standard git operations  
⚡ **Zero Configuration Conflicts**: All your existing git aliases, scripts, and workflows continue to work  

## ✨ Quick Example

See the power of sgit in action! Here's what happens when you replace `git` with `sgit`:

### Before (Traditional Git)
```bash
$ git commit -m "fix stuff" -a
[main a1b2c3d] fix stuff
 3 files changed, 47 insertions(+), 12 deletions(-)
```

### After (sgit with AI)
```bash
$ sgit commit -a
🤖 Staging all modified and deleted files...
🤖 Analyzing changes and generating commit message...

Generated commit message:
feat(auth): implement OAuth2 login with Google integration

Add Google OAuth2 authentication flow with token refresh mechanism.
Update login component to handle OAuth2 redirects and error states.
Include comprehensive test coverage for auth service methods.

? Use this commit message? (Y/n) y
[main e4f5g6h] feat(auth): implement OAuth2 login with Google integration
 3 files changed, 47 insertions(+), 12 deletions(-)
```

**The Result**: Instead of meaningless "fix stuff", you get professional conventional commits that clearly explain what changed and why! 🎉

## 🎯 Why Use sgit?

**Stop writing bad commit messages forever!** sgit transforms your development workflow:

### 🚀 **Instant Professional Commits**
- **Before**: `git commit -m "fix bug"`
- **After**: `feat(api): implement rate limiting with Redis backend`
- No more embarrassing commit history in code reviews!

### 🧠 **Smart File Management**  
- Automatically detects which files should be committed
- Skips temporary files, logs, and build artifacts
- Learns from your project structure

### 📈 **Better Code Reviews**
- Reviewers immediately understand what changed
- Conventional commits enable automated changelog generation
- Professional commit history impresses teammates and employers

### ⚡ **Zero Learning Curve**
- Use exactly like git: `sgit status`, `sgit push`, `sgit pull`
- All your existing git knowledge and muscle memory works
- Existing scripts and aliases continue working

### 🔥 **Real Examples from Users**
```bash
# Instead of this... 😞
git commit -m "updates"
git commit -m "fix"  
git commit -m "more changes"

# Get this automatically! ✨
feat(auth): add JWT token validation middleware
fix(db): resolve connection pooling memory leak  
docs(api): update authentication endpoint examples
```

## Features

- 🤖 **AI-powered commit messages**: Automatically generates comprehensive conventional commit messages using Solar LLM
- 🧠 **Smart file staging**: AI analyzes untracked files to decide what should be added to git  
- 📊 **Intelligent diff summaries**: AI explains what changed in your diffs
- 📈 **Log analysis**: AI analyzes commit history patterns and provides insights
- 🔀 **Merge assistance**: AI helps resolve conflicts and generates merge commit messages
- 🔄 **Full Git Compatibility**: Supports ALL git commands and options - just replace `git` with `sgit`
- ⚡ **Easy setup**: Simple configuration with your Upstage API key
- 🎛️ **Flexible options**: Interactive mode, manual override, and traditional git commit fallback
- 📦 **Single binary**: No dependencies, just download and run
- 🛡️ **Smart filtering**: Automatically skips binary files and large files

## 🚀 Try It Now!

**Ready to upgrade your git experience?** Get started in 30 seconds:

```bash
# 1. Install sgit (one command)
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash

# 2. Use it exactly like git
cd your-project
sgit status
sgit commit -a           # ✨ Stage files + AI-generated commit message

# 3. Enjoy professional commits!
```

**Want to see the difference immediately?** Try these commands in any git repository:
- `sgit log --ai-analysis` - Get insights about your development patterns  
- `sgit diff --ai-summary` - Understand what changed in plain English
- `sgit add --all-ai --dry-run-ai` - See which files AI recommends adding

**No API key?** No problem! Get your free key at [Upstage Console](https://console.upstage.ai/) - takes 2 minutes.

## Installation

sgit provides multiple installation methods to suit different preferences and platforms:

### 🚀 Quick Install (Recommended)

One-liner installation script that automatically detects your platform:

```bash
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash
```

### 🍺 Homebrew (macOS & Linux)

```bash
# Add the tap
brew tap hunkim/sgit

# Install sgit
brew install sgit

# Upgrade to latest version
brew upgrade sgit
```

### 📦 Package Managers

#### **Go Install**
If you have Go installed:
```bash
go install github.com/hunkim/sgit@latest
```

#### **Linux Package Managers**

**Debian/Ubuntu (APT):**
```bash
# Download .deb package from releases
wget https://github.com/hunkim/sgit/releases/latest/download/sgit_linux_amd64.deb
sudo dpkg -i sgit_linux_amd64.deb
```

**Red Hat/Fedora/CentOS (RPM):**
```bash
# Download .rpm package from releases  
wget https://github.com/hunkim/sgit/releases/latest/download/sgit_linux_amd64.rpm
sudo rpm -i sgit_linux_amd64.rpm
```

**Arch Linux (AUR):**
```bash
# Using yay
yay -S sgit

# Using paru
paru -S sgit
```

#### **Windows Package Managers**

**Chocolatey:**
```powershell
choco install sgit
```

**Scoop:**
```powershell
scoop bucket add hunkim https://github.com/hunkim/scoop-bucket
scoop install sgit
```

**Winget:**
```powershell
winget install hunkim.sgit
```

### 📥 Manual Download

Download pre-built binaries for your platform:

1. Go to [Releases](https://github.com/hunkim/sgit/releases/latest)
2. Download the appropriate binary:
   - **Linux:** `sgit_vX.X.X_linux_amd64.tar.gz`
   - **macOS:** `sgit_vX.X.X_darwin_amd64.tar.gz` (Intel) or `sgit_vX.X.X_darwin_arm64.tar.gz` (Apple Silicon)
   - **Windows:** `sgit_vX.X.X_windows_amd64.zip`
3. Extract and add to your PATH

### 🔨 Build from Source

```bash
git clone https://github.com/hunkim/sgit.git
cd sgit
go build -o sgit
sudo mv sgit /usr/local/bin/  # or add to your PATH
```

### 🐳 Docker

```bash
# Run sgit in a container
docker run --rm -it -v $(pwd):/workspace hunkim/sgit:latest

# Create an alias for convenience
echo 'alias sgit="docker run --rm -it -v \$(pwd):/workspace hunkim/sgit:latest"' >> ~/.bashrc
```

### 📋 Installation Verification

After installation, verify sgit is working:

```bash
# Check installation
sgit --help

# Verify git compatibility
sgit status

# Test AI features (requires API key setup)
sgit config
```

### 🔄 Updating sgit

**Homebrew:**
```bash
brew upgrade sgit
```

**Quick Install Script:**
```bash
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash
```

**Go Install:**
```bash
go install github.com/hunkim/sgit@latest
```

### ❌ Uninstallation

**Homebrew:**
```bash
brew uninstall sgit
brew untap hunkim/sgit
```

**Manual:**
```bash
# Remove binary
sudo rm /usr/local/bin/sgit  # or wherever you installed it

# Remove config (optional)
rm -rf ~/.config/sgit
```

---

## Setup

sgit will automatically prompt you to configure your API key when you first use any AI-powered feature:

```bash
# First time using sgit - automatic setup
sgit commit
# or
sgit add --all-ai
```

You can also manually configure at any time:

```bash
sgit config
```

This will prompt you for:
- **Upstage API Key**: Your API key from [Upstage Console](https://console.upstage.ai/)
- **Model Name**: Solar model to use (default: `solar-pro2-preview`)

Your API key is stored locally and securely in `~/.config/sgit/config.yaml`.

**Note**: Commands that don't require AI work without configuration - sgit is a full git replacement.

## Git Compatibility

**sgit is 100% compatible with git.** You can use it as a drop-in replacement:

```bash
# These work exactly the same as git (automatic passthrough)
sgit status
sgit branch
sgit push
sgit pull
sgit checkout main
sgit remote -v
sgit stash
sgit tag

# AI-enhanced commands (when available)
sgit commit                    # Uses AI-enhanced commit
sgit add --all-ai             # Uses AI-enhanced add
sgit diff --ai-summary        # Uses AI-enhanced diff
sgit log --ai-analysis        # Uses AI-enhanced log
sgit merge --ai-help          # Uses AI-enhanced merge

# Force standard git behavior for implemented commands
sgit git commit -m "message"  # Standard git commit
sgit git log --oneline -5     # Standard git log
sgit git diff HEAD~1          # Standard git diff

# All git flags work seamlessly
sgit commit --amend           # AI-enhanced with git flags
sgit add -p                   # Standard git (no AI flags)
```

**Key Points:**
- ✅ **Unimplemented commands**: Automatically pass through to git
- ✅ **Implemented commands**: Use AI-enhanced versions by default
- ✅ **All git flags**: Fully supported in both modes
- ✅ **Explicit passthrough**: Use `sgit git <command>` for standard git behavior

The AI features are **opt-in** and only activate when you explicitly use AI-specific flags or behavior.

## AI-Enhanced Commands

### 1. Smart Add (`sgit add`)

AI-powered file staging that analyzes untracked files:

```bash
# Analyze all untracked files and get AI recommendations
sgit add --all-ai

# Preview what would be added without actually adding
sgit add --all-ai --dry-run-ai

# Add files without AI confirmation (based on file type detection only)
sgit add --all-ai --force-ai

# Force AI analysis even for specific files
sgit add file1.js file2.py --ai

# Traditional behavior still works
sgit add file1.go file2.js  # No AI, just like git
```

### 2. Comprehensive Commits (`sgit commit`)

Enhanced AI-powered commit message generation with context awareness:

```bash
# Default: AI generates comprehensive message considering:
# - Git diff, branch name, recent commits, files changed
# - Opens editor with AI-generated message pre-filled
sgit commit

# Stage all modified files and generate AI commit message
sgit commit -a

# AI commit with amend (edit previous commit with AI message)
sgit commit --amend

# AI commit with verbose output
sgit commit -v

# AI commit with signoff
sgit commit --signoff

# Skip editor, use AI message with confirmation
sgit commit --skip-editor

# Stage all files and skip editor (quick workflow)
sgit commit -a --skip-editor

# Interactive terminal editing of AI message
sgit commit -i

# Traditional git behavior (NO AI) - explicit opt-out
sgit commit -m "manual message"  # Manual message bypasses AI
sgit commit --no-ai             # Explicit no-AI flag
```

### 3. Intelligent Diff (`sgit diff`)

Enhanced diff with AI-powered summaries:

```bash
# Show diff with AI-powered summary
sgit diff --ai-summary

# Works with all git diff options
sgit diff --cached --ai-summary
sgit diff HEAD~1 --ai-summary
sgit diff --stat --ai-summary

# Traditional git diff (no AI)
sgit diff
sgit diff --cached
sgit diff HEAD~1
```

### 4. Log Analysis (`sgit log`)

Commit history analysis with AI insights:

```bash
# Show log with AI-powered analysis
sgit log --ai-analysis

# Analyze specific timeframe
sgit log --ai-analysis --ai-timeframe "last 2 weeks"

# Works with all git log options
sgit log --oneline --ai-analysis
sgit log --since="1 week ago" --ai-analysis
sgit log --author="John" --ai-analysis

# Traditional git log (no AI)
sgit log
sgit log --oneline
sgit log --graph
```

### 5. Merge Assistance (`sgit merge`)

AI-powered merge conflict resolution and commit messages:

```bash
# Merge with AI conflict assistance
sgit merge feature-branch --ai-help

# Merge with AI-generated merge commit message
sgit merge feature-branch --ai-message

# Both AI assistance and message
sgit merge feature-branch --ai-help --ai-message

# Traditional git merge (no AI)
sgit merge feature-branch
sgit merge --no-ff feature-branch
```

## Complete Workflow Examples

```bash
# Full AI-enhanced workflow
sgit add --all-ai              # Smart file staging
sgit commit                    # Comprehensive AI commit message
sgit diff --ai-summary         # Review changes with AI
sgit log --ai-analysis         # Analyze development patterns

# Quick AI workflow
sgit add --all-ai --force-ai   # Smart add without confirmation
sgit commit --skip-editor      # AI message with confirmation

# Ultra-streamlined workflow (most common)
sgit commit -a --skip-editor   # Stage all modified files + AI commit in one command

# Development analysis workflow
sgit log --ai-analysis --ai-timeframe "last sprint"
sgit diff --cached --ai-summary
sgit merge feature --ai-help --ai-message

# Traditional workflow (no AI at all)
sgit add .                     # Standard git add
sgit commit -m "message"       # Standard git commit (manual message)
sgit commit --no-ai            # Standard git commit (explicit no-AI)
sgit merge feature-branch      # Standard git merge
```

## How AI Features Work

### Smart Add (`--all-ai`)
1. **File Discovery**: Finds all untracked files using `git ls-files --others --exclude-standard`
2. **Binary Detection**: Skips files based on extension and content analysis
3. **Size Filtering**: Skips files larger than 1MB
4. **AI Analysis**: Sends file content to Solar LLM for decision
5. **User Review**: Shows recommendations with reasons
6. **Execution**: Adds approved files to staging area

### Comprehensive Commits (default behavior)
1. **Context Gathering**: Collects git diff, branch name, recent commits, file list
2. **AI Processing**: Sends comprehensive context to Solar LLM
3. **Message Generation**: Solar LLM generates detailed conventional commit message
4. **User Review**: Opens editor with AI message pre-filled for review/editing
5. **Git Commit**: Executes `git commit` with final message

### Diff Summaries (`--ai-summary`)
1. **Diff Analysis**: Captures git diff output
2. **AI Processing**: Sends diff to Solar LLM for analysis
3. **Summary Generation**: AI provides structured summary of changes
4. **Display**: Shows both original diff and AI summary

### Log Analysis (`--ai-analysis`)
1. **Log Collection**: Gathers git log output (last 20 commits by default)
2. **Pattern Analysis**: AI analyzes commit patterns, contributors, changes
3. **Insight Generation**: Provides development insights and recommendations
4. **Display**: Shows both log and AI analysis

### Merge Assistance (`--ai-help`, `--ai-message`)
1. **Conflict Detection**: Identifies merge conflicts automatically
2. **AI Guidance**: Provides resolution strategies and risk assessment
3. **Message Generation**: Creates comprehensive merge commit messages
4. **User Support**: Guides through conflict resolution process

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
- `perf:` - Performance improvements
- `ci:` - CI/CD changes
- `build:` - Build system changes

## Editor Configuration

sgit uses the same editor configuration as git, in this order of preference:

1. `GIT_EDITOR` environment variable
2. `git config core.editor`
3. `VISUAL` environment variable  
4. `EDITOR` environment variable
5. Falls back to `nano`, `vim`, or `vi`

To set your preferred editor:
```bash
# For sgit and git
git config --global core.editor "code --wait"  # VS Code
git config --global core.editor "vim"          # Vim
git config --global core.editor "nano"         # Nano

# Or set environment variable
export EDITOR="code --wait"
```

## Configuration

Configuration file location: `~/.config/sgit/config.yaml`

```yaml
upstage_api_key: "up_****************************"
upstage_model_name: "solar-pro2-preview"
```

## API Key Setup

1. Sign up at [Upstage Console](https://console.upstage.ai/)
2. Navigate to API Keys section in the console
3. Create a new API key for sgit
4. Run `sgit config` (or use any AI command for automatic setup)
5. Enter your API key when prompted

Your API key will be stored locally and securely on your machine.

## Requirements

- Git installed and configured
- Upstage API key (only for AI features)
- Internet connection (only for AI features)

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

**Note**: This tool requires an Upstage API key only for AI features. All standard git functionality works without any configuration or API key. 