#!/bin/bash

# Build release binaries for multiple platforms
# This script should be run from the project root

set -e

VERSION=${1:-"dev"}
OUTPUT_DIR="dist"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Create output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

print_status "Building sgit version $VERSION for multiple platforms..."

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    output_name="sgit"
    if [ "$GOOS" = "windows" ]; then
        output_name="sgit.exe"
    fi
    
    archive_name="sgit_${VERSION}_${GOOS}_${GOARCH}"
    
    print_status "Building for $GOOS/$GOARCH..."
    
    # Build binary
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="-s -w -X main.version=$VERSION" \
        -o "$OUTPUT_DIR/$output_name" \
        .
    
    # Create archive
    if [ "$GOOS" = "windows" ]; then
        # Create ZIP for Windows
        cd "$OUTPUT_DIR"
        zip -q "${archive_name}.zip" "$output_name"
        rm "$output_name"
        cd ..
        print_success "Created ${archive_name}.zip"
    else
        # Create tar.gz for Unix-like systems
        cd "$OUTPUT_DIR"
        tar -czf "${archive_name}.tar.gz" "$output_name"
        rm "$output_name"
        cd ..
        print_success "Created ${archive_name}.tar.gz"
    fi
done

print_success "Build complete! Release files are in $OUTPUT_DIR/"
echo ""
print_status "Files created:"
ls -la "$OUTPUT_DIR"/

# Generate checksums
print_status "Generating checksums..."
cd "$OUTPUT_DIR"
sha256sum * > checksums.txt
cd ..
print_success "Checksums saved to $OUTPUT_DIR/checksums.txt" 