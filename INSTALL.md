# Installation Guide for sgit Development

This guide will help you set up the development environment and build sgit from source.

## Prerequisites

### 1. Install Go (Required)

**macOS:**
```bash
# Using Homebrew (recommended)
brew install go

# Or download from https://golang.org/dl/
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL/Fedora
sudo dnf install golang
# or
sudo yum install golang

# Or download from https://golang.org/dl/
```

**Windows:**
- Download the installer from https://golang.org/dl/
- Run the installer and follow the instructions

### 2. Verify Go Installation

```bash
go version
```
You should see output like: `go version go1.21.0 darwin/amd64`

### 3. Install Git (if not already installed)

**macOS:**
```bash
brew install git
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install git

# CentOS/RHEL/Fedora
sudo dnf install git
```

**Windows:**
- Download from https://git-scm.com/

## Building sgit

### 1. Clone the Repository

```bash
git clone https://github.com/hunkim/sgit.git
cd sgit
```

### 2. Download Dependencies

```bash
go mod tidy
```

### 3. Build the Project

**Quick build:**
```bash
make build
```

**Or manually:**
```bash
go build -o build/sgit .
```

### 4. Test the Build

```bash
./build/sgit --help
```

### 5. Install Locally (Optional)

```bash
make install
# or
go install .
```

This will install sgit to your `$GOPATH/bin` directory.

## Cross-Platform Build

To build for all supported platforms:

```bash
make build-all
```

This creates binaries for:
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (AMD64)

## Creating Release Packages

```bash
make release
```

This creates compressed archives in the `build/` directory ready for distribution.

## Development Workflow

### 1. Make changes to the code

### 2. Test locally
```bash
make quick
```

### 3. Run tests (when tests are added)
```bash
make test
```

### 4. Build and test
```bash
make build
./build/sgit config  # Setup your API key
./build/sgit commit  # Test the functionality
```

## Troubleshooting

### Go Command Not Found
- Ensure Go is properly installed and in your PATH
- Restart your terminal after installation
- Check `echo $PATH` includes Go's bin directory

### Module Download Issues
- Ensure you have internet connectivity
- Try `go clean -modcache` and then `go mod tidy`

### Build Errors
- Make sure you're using Go 1.21 or later
- Check that all source files are present
- Try `make clean` and rebuild

### Permission Issues (macOS/Linux)
```bash
chmod +x build/sgit
```

## API Key Setup

After building, configure your Upstage API key:

1. Get your API key from [Upstage Console](https://console.upstage.ai/)
2. Run: `./build/sgit config`
3. Enter your API key when prompted

## Development Tips

- Use `make quick` for rapid development cycles
- Configuration is stored in `~/.config/sgit/config.yaml`
- Add `build/` directory to your PATH for easier testing
- Use `-h` or `--help` with any command to see options

## Next Steps

Once built and configured:

1. Navigate to a git repository
2. Make some changes and stage them: `git add .`
3. Test sgit: `./build/sgit commit`
4. Enjoy AI-generated commit messages!

For usage instructions, see the main [README.md](README.md). 