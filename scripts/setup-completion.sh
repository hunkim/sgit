#!/bin/bash

# Setup script for sgit completion
# This script sets up tab completion for sgit commands

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SGIT_BINARY="${SCRIPT_DIR}/../sgit"

# Check if sgit binary exists
if [ ! -f "$SGIT_BINARY" ]; then
    echo "❌ sgit binary not found at $SGIT_BINARY"
    echo "Please build sgit first: go build -o sgit"
    exit 1
fi

# Detect shell
if [ -n "$ZSH_VERSION" ]; then
    SHELL_TYPE="zsh"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_TYPE="bash"
else
    # Try to detect from $SHELL variable
    case "$SHELL" in
        */zsh)
            SHELL_TYPE="zsh"
            ;;
        */bash)
            SHELL_TYPE="bash"
            ;;
        */fish)
            SHELL_TYPE="fish"
            ;;
        *)
            echo "❓ Unable to detect shell type. Please specify: bash, zsh, or fish"
            echo "Usage: $0 [bash|zsh|fish]"
            exit 1
            ;;
    esac
fi

# Allow override from command line
if [ $# -gt 0 ]; then
    SHELL_TYPE="$1"
fi

echo "🚀 Setting up sgit completion for $SHELL_TYPE..."

case "$SHELL_TYPE" in
    "bash")
        COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
        COMPLETION_FILE="$COMPLETION_DIR/sgit"
        RC_FILE="$HOME/.bashrc"
        
        # Create completion directory
        mkdir -p "$COMPLETION_DIR"
        
        # Generate completion
        "$SGIT_BINARY" completion bash > "$COMPLETION_FILE"
        
        echo "✅ Bash completion installed to $COMPLETION_FILE"
        echo "🔄 Please restart your shell or run: source ~/.bashrc"
        ;;
        
    "zsh")
        COMPLETION_DIR="$HOME/.local/share/zsh/site-functions"
        COMPLETION_FILE="$COMPLETION_DIR/_sgit"
        RC_FILE="$HOME/.zshrc"
        
        # Create completion directory
        mkdir -p "$COMPLETION_DIR"
        
        # Generate completion
        "$SGIT_BINARY" completion zsh > "$COMPLETION_FILE"
        
        # Add completion directory to fpath if not already there
        if ! grep -q "fpath=.*$COMPLETION_DIR" "$RC_FILE" 2>/dev/null; then
            echo "" >> "$RC_FILE"
            echo "# sgit completion" >> "$RC_FILE"
            echo "fpath=(\"$COMPLETION_DIR\" \$fpath)" >> "$RC_FILE"
            echo "autoload -U compinit && compinit" >> "$RC_FILE"
        fi
        
        echo "✅ Zsh completion installed to $COMPLETION_FILE"
        echo "✅ Added completion directory to $RC_FILE"
        echo "🔄 Please restart your shell or run: source ~/.zshrc"
        ;;
        
    "fish")
        COMPLETION_DIR="$HOME/.config/fish/completions"
        COMPLETION_FILE="$COMPLETION_DIR/sgit.fish"
        
        # Create completion directory
        mkdir -p "$COMPLETION_DIR"
        
        # Generate completion
        "$SGIT_BINARY" completion fish > "$COMPLETION_FILE"
        
        echo "✅ Fish completion installed to $COMPLETION_FILE"
        echo "🔄 Completion will be available in new fish sessions"
        ;;
        
    *)
        echo "❌ Unsupported shell: $SHELL_TYPE"
        echo "Supported shells: bash, zsh, fish"
        exit 1
        ;;
esac

echo ""
echo "🎉 sgit completion setup complete!"
echo ""
echo "📝 Available sgit commands:"
echo "  • sgit add       - Add files with AI analysis"
echo "  • sgit commit    - Commit with AI-generated messages"  
echo "  • sgit diff      - Show changes with AI summary"
echo "  • sgit log       - Show logs with AI analysis"
echo "  • sgit merge     - Merge with AI assistance"
echo "  • sgit config    - Configure settings"
echo ""
echo "🌍 Language support:"
echo "  Use --lang flag: sgit --lang ko commit -a"
echo "  Available: en, ko, ja, zh, es, fr, de"
echo ""
echo "💡 Test completion:"
echo "  Type 'sgit <TAB>' to see available commands"
echo "  Type 'sgit --lang <TAB>' to see language options" 