#!/bin/bash

# sgit installation script
# Usage: curl -fsSL https://raw.githubusercontent.com/hunkim/sgit/main/scripts/install.sh | bash
# Options:
#   --no-completion    Skip tab completion setup

set -e

# Parse command line arguments
SKIP_COMPLETION=false
for arg in "$@"; do
    case $arg in
        --no-completion)
            SKIP_COMPLETION=true
            shift
            ;;
        --help|-h)
            echo "sgit installation script"
            echo ""
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --no-completion    Skip tab completion setup"
            echo "  --help, -h         Show this help message"
            echo ""
            exit 0
            ;;
        *)
            # Unknown option
            echo "Unknown option: $arg"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $OS in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        windows*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}_${ARCH}"
}

# Get latest release version
get_latest_version() {
    print_status "Fetching latest release information..."
    
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s https://api.github.com/repos/hunkim/sgit/releases/latest | grep -o '"tag_name": "[^"]*' | grep -o '[^"]*$')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- https://api.github.com/repos/hunkim/sgit/releases/latest | grep -o '"tag_name": "[^"]*' | grep -o '[^"]*$')
    else
        print_error "curl or wget is required to download sgit"
        exit 1
    fi
    
    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi
    
    print_success "Latest version: $VERSION"
}

# Download and install binary
install_binary() {
    print_status "Downloading sgit for $PLATFORM..."
    
    DOWNLOAD_URL="https://github.com/hunkim/sgit/releases/download/${VERSION}/sgit_${VERSION}_${PLATFORM}.tar.gz"
    TEMP_DIR=$(mktemp -d)
    
    if command -v curl >/dev/null 2>&1; then
        curl -L "$DOWNLOAD_URL" -o "$TEMP_DIR/sgit.tar.gz"
    else
        wget "$DOWNLOAD_URL" -O "$TEMP_DIR/sgit.tar.gz"
    fi
    
    if [ $? -ne 0 ]; then
        print_error "Failed to download sgit"
        rm -rf "$TEMP_DIR"
        exit 1
    fi
    
    print_status "Extracting binary..."
    tar -xzf "$TEMP_DIR/sgit.tar.gz" -C "$TEMP_DIR"
    
    # Determine installation directory
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    elif [ -w "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    else
        print_warning "Neither /usr/local/bin nor ~/.local/bin is writable"
        print_status "Please run with sudo or create ~/.local/bin and add it to PATH"
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
    
    print_status "Installing sgit to $INSTALL_DIR..."
    
    if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w "/usr/local/bin" ]; then
        sudo cp "$TEMP_DIR/sgit" "$INSTALL_DIR/sgit"
        sudo chmod +x "$INSTALL_DIR/sgit"
    else
        cp "$TEMP_DIR/sgit" "$INSTALL_DIR/sgit"
        chmod +x "$INSTALL_DIR/sgit"
    fi
    
    rm -rf "$TEMP_DIR"
    
    print_success "sgit installed successfully!"
}

# Check if git is installed
check_git() {
    if ! command -v git >/dev/null 2>&1; then
        print_error "git is required but not installed. Please install git first."
        case $OS in
            linux)
                echo "  Ubuntu/Debian: sudo apt install git"
                echo "  CentOS/RHEL: sudo yum install git"
                echo "  Fedora: sudo dnf install git"
                ;;
            darwin)
                echo "  macOS: brew install git (or install Xcode Command Line Tools)"
                ;;
        esac
        exit 1
    fi
}

# Check PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        print_warning "$INSTALL_DIR is not in your PATH"
        print_status "Add the following to your shell profile (.bashrc, .zshrc, etc.):"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
        print_status "Or run: echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.$(basename $SHELL)rc"
    fi
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    if command -v sgit >/dev/null 2>&1; then
        VERSION_OUTPUT=$(sgit --version 2>/dev/null || sgit --help | head -1)
        print_success "sgit is working: $VERSION_OUTPUT"
        return 0
    else
        print_warning "sgit command not found in PATH"
        if [ -f "$INSTALL_DIR/sgit" ]; then
            print_status "sgit is installed at $INSTALL_DIR/sgit"
            check_path
            return 1
        else
            print_error "Installation verification failed"
            return 1
        fi
    fi
}

# Setup tab completion
setup_completion() {
    print_status "Setting up tab completion..."
    
    # Skip completion setup if sgit is not in PATH
    if ! command -v sgit >/dev/null 2>&1; then
        print_warning "Skipping completion setup - sgit not found in PATH"
        print_status "Run 'sgit completion --help' after adding sgit to PATH to set up completion manually"
        return
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
                print_warning "Unable to detect shell type for completion setup"
                print_status "Run 'sgit completion --help' to set up completion manually"
                return
                ;;
        esac
    fi
    
    print_status "Detected shell: $SHELL_TYPE"
    
    case "$SHELL_TYPE" in
        "bash")
            COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
            COMPLETION_FILE="$COMPLETION_DIR/sgit"
            RC_FILE="$HOME/.bashrc"
            
            # Create completion directory
            mkdir -p "$COMPLETION_DIR"
            
            # Generate completion
            if sgit completion bash > "$COMPLETION_FILE" 2>/dev/null; then
                print_success "Bash completion installed to $COMPLETION_FILE"
                print_status "Completion will be available in new bash sessions"
            else
                print_warning "Failed to generate bash completion"
            fi
            ;;
            
        "zsh")
            COMPLETION_DIR="$HOME/.local/share/zsh/site-functions"
            COMPLETION_FILE="$COMPLETION_DIR/_sgit"
            RC_FILE="$HOME/.zshrc"
            
            # Create completion directory
            mkdir -p "$COMPLETION_DIR"
            
            # Generate completion
            if sgit completion zsh > "$COMPLETION_FILE" 2>/dev/null; then
                print_success "Zsh completion installed to $COMPLETION_FILE"
                
                # Add completion directory to fpath if not already there
                if [ -f "$RC_FILE" ] && ! grep -q "fpath=.*$COMPLETION_DIR" "$RC_FILE" 2>/dev/null; then
                    echo "" >> "$RC_FILE"
                    echo "# sgit completion" >> "$RC_FILE"
                    echo "fpath=(\"$COMPLETION_DIR\" \$fpath)" >> "$RC_FILE"
                    echo "autoload -U compinit && compinit" >> "$RC_FILE"
                    print_success "Added completion directory to $RC_FILE"
                elif [ ! -f "$RC_FILE" ]; then
                    echo "# sgit completion" > "$RC_FILE"
                    echo "fpath=(\"$COMPLETION_DIR\" \$fpath)" >> "$RC_FILE"
                    echo "autoload -U compinit && compinit" >> "$RC_FILE"
                    print_success "Created $RC_FILE with completion setup"
                fi
                
                print_status "Completion will be available in new zsh sessions"
            else
                print_warning "Failed to generate zsh completion"
            fi
            ;;
            
        "fish")
            COMPLETION_DIR="$HOME/.config/fish/completions"
            COMPLETION_FILE="$COMPLETION_DIR/sgit.fish"
            
            # Create completion directory
            mkdir -p "$COMPLETION_DIR"
            
            # Generate completion
            if sgit completion fish > "$COMPLETION_FILE" 2>/dev/null; then
                print_success "Fish completion installed to $COMPLETION_FILE"
                print_status "Completion will be available in new fish sessions"
            else
                print_warning "Failed to generate fish completion"
            fi
            ;;
            
        *)
            print_warning "Unsupported shell for automatic completion setup: $SHELL_TYPE"
            print_status "Run 'sgit completion --help' to set up completion manually"
            ;;
    esac
}

# Main installation flow
main() {
    echo "Installing sgit - Solar LLM-powered git wrapper"
    echo "=============================================="
    echo ""
    
    check_git
    detect_platform
    get_latest_version
    install_binary
    
    # Try to set up completion if installation was successful
    if verify_installation; then
        if [ "$SKIP_COMPLETION" = false ]; then
            setup_completion
        else
            print_status "Skipping completion setup (--no-completion flag used)"
            print_status "Run 'sgit completion --help' to set up completion manually later"
        fi
    fi
    
    echo ""
    print_success "Installation complete!"
    echo ""
    print_status "🎉 sgit is ready to use!"
    echo ""
    print_status "✨ What's included:"
    echo "  • AI-powered commit messages"
    echo "  • Smart diff summaries"  
    echo "  • Intelligent log analysis"
    echo "  • Multi-language support (en, ko, ja, zh, es, fr, de)"
    echo "  • Tab completion for commands and flags"
    echo ""
    print_status "🚀 Next steps:"
    echo "  1. Run 'sgit config' to set up your Upstage API key"
    echo "  2. Try 'sgit commit -a' for AI-generated commit messages"
    echo "  3. Use 'sgit --lang ko commit' for Korean responses"
    echo "  4. Press TAB after 'sgit' to see available commands"
    echo ""
    print_status "🔗 Get your free API key at: https://console.upstage.ai/"
    echo ""
    print_status "💡 Tips:"
    echo "  • All git commands work: 'sgit status', 'sgit push', etc."
    echo "  • Use 'sgit --help' to see all features"
    echo "  • Restart your shell to enable tab completion"
    echo ""
}

# Run main function
main "$@" 