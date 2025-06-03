# 🚀 sgit Tab Completion

sgit supports intelligent tab completion for commands, flags, and language codes, just like git!

## ⚡ Quick Setup (Automatic)

**Tab completion is now included by default!** 🎉

When you install sgit using the main installation script, tab completion is automatically set up:

```bash
# This now includes tab completion setup automatically
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash
```

**What happens automatically:**
- ✅ Detects your shell (bash/zsh/fish)
- ✅ Installs completion in the correct location
- ✅ Configures your shell profile
- ✅ No additional steps needed!

**Skip completion if you don't want it:**
```bash
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash -s -- --no-completion
```

## 🛠 Manual Setup (Optional)

If you installed sgit through other methods (Homebrew, go install, etc.), you can set up completion manually:

**Quick manual setup:**
```bash
# Use the standalone completion script
./scripts/setup-completion.sh
```

**Or generate manually for your shell:**
```bash
# For bash
sgit completion bash > ~/.local/share/bash-completion/completions/sgit

# For zsh
mkdir -p ~/.local/share/zsh/site-functions
sgit completion zsh > ~/.local/share/zsh/site-functions/_sgit
echo 'fpath=(~/.local/share/zsh/site-functions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc

# For fish
sgit completion fish > ~/.config/fish/completions/sgit.fish
```

## 🎯 What Gets Completed

### Commands
- `sgit <TAB>` → Shows all available commands
  - `add` - Add files with AI analysis
  - `commit` - Commit with AI-generated messages
  - `diff` - Show changes with AI summary
  - `log` - Show logs with AI analysis
  - `merge` - Merge with AI assistance
  - `config` - Configure settings
  - `completion` - Generate completion scripts
  - `help` - Show help
  - `version` - Show version

### Global Flags
- `sgit --<TAB>` → Shows global flags
  - `--config` - Config file path
  - `--lang` - Language for AI responses
  - `--help` - Show help
  - `--version` - Show version

### Language Codes
- `sgit --lang <TAB>` → Shows available languages
  - `en` - English
  - `ko` - Korean (한국어)
  - `ja` - Japanese (日本語)
  - `zh` - Chinese (中文)
  - `es` - Spanish (Español)
  - `fr` - French (Français)
  - `de` - German (Deutsch)

### Command-Specific Flags
Each command has its own set of completable flags:

**Commit flags:**
- `sgit commit --<TAB>` → Shows commit-specific flags
  - `--all` / `-a` - Stage all modified files
  - `--interactive` / `-i` - Interactive mode
  - `--message` / `-m` - Commit message
  - `--amend` - Amend last commit
  - `--signoff` - Add signed-off-by line
  - And many more...

## 💡 Usage Examples

```bash
# Tab completion for commands
sgit <TAB>
# Shows: add commit config diff help log merge version completion

# Tab completion for language
sgit --lang <TAB>
# Shows: en ko ja zh es fr de

# Tab completion for commit flags
sgit commit --<TAB>
# Shows: --all --interactive --message --amend --signoff --no-ai --skip-editor ...

# Combined usage with completion
sgit --lang ko commit -a -i
#      ↑ tab complete   ↑ tab complete flags
```

## 🔧 Troubleshooting

### Completion not working after installation?

1. **Restart your shell:**
   ```bash
   # Open a new terminal or source your profile
   source ~/.bashrc    # For bash
   source ~/.zshrc     # For zsh
   ```

2. **Check if completion files exist:**
   ```bash
   # For bash
   ls ~/.local/share/bash-completion/completions/sgit
   
   # For zsh  
   ls ~/.local/share/zsh/site-functions/_sgit
   
   # For fish
   ls ~/.config/fish/completions/sgit.fish
   ```

3. **Manually reinstall completion:**
   ```bash
   ./scripts/setup-completion.sh
   ```

4. **Check if sgit is in PATH:**
   ```bash
   which sgit
   sgit --help
   ```

### Different installation methods

| Installation Method | Completion Included? | Manual Setup |
|-------------------|---------------------|--------------|
| **Quick Install Script** | ✅ **Automatic** | Not needed |
| Homebrew | ❌ Manual | Run `./scripts/setup-completion.sh` |
| Go Install | ❌ Manual | Run `./scripts/setup-completion.sh` |
| Package Managers | ❌ Manual | Run `./scripts/setup-completion.sh` |
| Build from Source | ❌ Manual | Run `./scripts/setup-completion.sh` |

## 🎉 Features

✅ **Smart command completion** - Only shows sgit native commands  
✅ **Language code completion** - Tab complete supported language codes  
✅ **Flag completion** - Complete command-specific flags  
✅ **Multi-shell support** - Works with bash, zsh, and fish  
✅ **Auto-detection** - Automatically detects your shell  
✅ **Easy setup** - Included in default installation  
✅ **Optional** - Can be skipped with `--no-completion` flag

## 🆚 Comparison with Git

| Feature | git | sgit |
|---------|-----|------|
| Command completion | ✅ All git commands | ✅ sgit native commands only |
| Flag completion | ✅ | ✅ |
| Language completion | ❌ | ✅ `--lang <TAB>` |
| AI flag completion | ❌ | ✅ `--no-ai`, `--skip-editor`, etc. |
| Setup complexity | Complex | ✅ **Automatic in install** |

## 📝 Notes

- Completion is **automatically included** in the main installation script
- Only includes **sgit native commands**, not all git commands  
- This is intentional to avoid overwhelming users with too many options
- For git-specific commands, use `git <command>` directly
- Language completion provides both codes and full names for clarity
- The completion respects the same command structure as git for familiarity

---

**🎯 Bottom Line**: Tab completion now works out-of-the-box with zero additional setup required! 