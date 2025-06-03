# sgit - AI-Powered Git That Writes Perfect Commits

> **Never write a bad commit message again.** sgit uses AI to generate professional conventional commits automatically while maintaining 100% git compatibility.

## 🎬 See the Magic

### Your Current Git Workflow 😞
```bash
$ git add .
$ git commit -m "fix stuff"
$ git commit -m "updates"  
$ git commit -m "more changes"
```

### With sgit ✨
```bash
$ sgit commit -a
🤖 Analyzing changes...

feat(auth): implement OAuth2 login with Google integration

Add secure OAuth2 authentication flow with token refresh mechanism.
Update login component to handle redirects and error states gracefully.
Include comprehensive test coverage for authentication service.

✨ Resolves #123, improves security architecture

? Use this commit? (Y/n) █
```

**Result**: Transform meaningless commits into professional documentation that impresses reviewers and future developers! 🚀

---

## ⚡ Quick Start (30 seconds)

```bash
# 1. Install
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash

# 2. Use like git, but better
sgit commit -a        # 🤖 AI writes your commit message
sgit diff            # 📊 AI explains what changed  
sgit log             # 📈 AI analyzes your patterns
```

**That's it!** All your existing git commands work exactly the same: `sgit status`, `sgit push`, `sgit pull`, etc.

---

## 🌟 Why Developers Love sgit

- **🎯 Perfect Commits**: Conventional commits with context, not "fix stuff"
- **⚡ Zero Learning**: Drop-in git replacement - use your existing knowledge
- **🌍 Multi-Language**: AI responds in 7+ languages (`--lang ko` for Korean!)
- **🔄 100% Compatible**: All git commands work - scripts, aliases, everything
- **⌨️ Smart Completion**: Tab completion for commands and language codes
- **🛡️ Privacy First**: Your code stays local, only diffs sent for analysis

---

## 🚀 Features That Transform Your Workflow

| Traditional Git | sgit Enhancement |
|----------------|------------------|
| `git commit -m "fix"` | **AI generates**: `fix(api): resolve rate limiting edge case in Redis backend` |
| `git diff` | **AI explains**: "Refactored authentication middleware to use JWT validation" |
| `git log` | **AI insights**: "Focus on security improvements, 3 bug fixes this week" |
| Manual file staging | **AI recommends**: "Add these 4 source files, skip temp files" |

---

## 🛠️ Installation Options

**🚀 Recommended (includes tab completion):**
```bash
curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash
```

**🍺 Homebrew:**
```bash
brew tap hunkim/sgit && brew install sgit
```

**📦 Go Install:**
```bash
go install github.com/hunkim/sgit@latest
```

**🐳 Docker:**
```bash
docker run --rm -it -v $(pwd):/workspace hunkim/sgit:latest
```

---

## 🔧 Configuration (2 minutes)

```bash
sgit config  # Set up your free Upstage API key
```

Get your API key at [console.upstage.ai](https://console.upstage.ai/) (free tier available).

---

## 🎯 Core Commands

### Smart Commits
```bash
sgit commit              # AI writes commit message  
sgit commit -a           # Stage all + AI commit
sgit commit -a --lang ko # Korean AI responses
```

### Intelligent Analysis  
```bash
sgit diff                # AI explains changes
sgit log                 # AI analyzes patterns
sgit add --all-ai        # AI recommends files to stage
```

### Traditional Git (unchanged)
```bash
sgit status              # Same as git status
sgit push                # Same as git push  
sgit pull                # Same as git pull
```

---

## 🌍 Multi-Language Support

```bash
sgit --lang ko commit    # Korean: "기능: 사용자 인증 시스템 구현"
sgit --lang ja diff      # Japanese: "変更内容の分析..."  
sgit --lang es log       # Spanish: "Análisis de patrones..."
```

**Supported**: English, Korean, Japanese, Chinese, Spanish, French, German

---

## 🔒 Privacy & Security

- ✅ **Local First**: Your code stays on your machine
- ✅ **Diff Only**: Only git diffs sent for commit message generation
- ✅ **No Storage**: Upstage doesn't store your code or diffs
- ✅ **Open Source**: Full transparency, audit the code yourself

---

## 📚 Examples

### Before sgit 😱
```
git log --oneline -5
a1b2c3d fix stuff
e4f5g6h updates  
h7i8j9k more changes
k1l2m3n bug fix
n4o5p6q refactor
```

### After sgit 🎉  
```
sgit log --oneline -5
a1b2c3d feat(auth): implement OAuth2 with Google integration
e4f5g6h fix(db): resolve connection pooling memory leak
h7i8j9k docs(api): add comprehensive endpoint documentation  
k1l2m3n perf(cache): optimize Redis operations for 40% speed gain
n4o5p6q refactor(ui): modernize component architecture with hooks
```

**Night and day difference!** 🌟

---

## 🆚 vs Other Tools

| Feature | sgit | Conventional Commits | Other AI Tools |
|---------|------|---------------------|----------------|
| **Zero Learning Curve** | ✅ | ❌ | ❌ |
| **Full Git Compatibility** | ✅ | ✅ | ❌ |  
| **AI Commit Messages** | ✅ | ❌ | ✅ |
| **Multi-Language** | ✅ | ❌ | ❌ |
| **Privacy Focused** | ✅ | ✅ | ❌ |
| **Tab Completion** | ✅ | ❌ | ❌ |

---

## 🤝 Contributing

```bash
git clone https://github.com/hunkim/sgit.git
cd sgit
go build
sgit commit -a  # Use sgit to contribute to sgit! 😄
```

---

## ⭐ Love sgit?

- 🌟 **Star this repo** if sgit saves you time
- 🐛 **Report issues** at [github.com/hunkim/sgit/issues](https://github.com/hunkim/sgit/issues)  
- 💡 **Request features** - we're always improving
- 🗣️ **Spread the word** - help other developers write better commits

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

**Created with ❤️ by developers, for developers.**

---

> **Ready to transform your git workflow?** Try sgit today and never write "fix stuff" again! 🚀
>
> ```bash
> curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash
> ``` 