# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2024-12-19

### Added
- Initial release of sgit - Solar LLM-powered git wrapper
- AI-powered commit message generation with streaming display
- Smart file staging with AI analysis (`sgit add --all-ai`)
- AI-enhanced diff summaries (`sgit diff --ai-summary`)
- Git log analysis with AI insights (`sgit log --ai-analysis`)
- Merge assistance and conflict resolution guidance (`sgit merge --ai-help`)
- Full git compatibility with automatic command passthrough
- Cross-platform support (Linux, macOS, Windows)
- Multiple installation methods (Homebrew, curl script, Go install)
- Beautiful loading animations with terminal-specific spinners
- Comprehensive context awareness (branch, recent commits, file contents)
- Configuration management for Upstage Solar LLM API
- Interactive and non-interactive modes for all AI features
- Complete documentation and examples

### Features
- **AI Commands**: 5 AI-enhanced git commands (add, commit, diff, log, merge)
- **Git Compatibility**: 100% compatible with existing git workflows
- **Installation**: Multiple package manager support and one-liner install
- **User Experience**: Progressive loading, clear feedback, professional UI
- **Configuration**: Automatic API key setup and secure storage
- **Documentation**: Comprehensive README with usage examples

### Technical
- Built with Go for single-binary deployment
- Cobra CLI framework for extensible command structure
- Upstage Solar LLM integration with streaming support
- Cross-platform build automation with GitHub Actions
- Homebrew formula and installation scripts
- Thread-safe spinner animations with terminal detection

## [Unreleased]
- Package manager submissions (APT, RPM, AUR, Chocolatey, etc.)
- Docker image support
- Additional AI features and improvements
- Performance optimizations 