#!/bin/bash
# APT Repository Setup Script for AIask
# Run this on a Linux machine to create and deploy the APT repository
#
# Prerequisites:
# - Go 1.23+ installed
# - git configured with push access to Hermithic/aiask
# - dpkg-deb available (standard on Debian/Ubuntu)
#
# Usage: ./setup-apt-repo.sh

set -e

VERSION="2.0.1"
REPO_NAME="aiask"
GITHUB_USER="Hermithic"
REPO_URL="https://github.com/${GITHUB_USER}/${REPO_NAME}.git"

echo "🚀 AIask APT Repository Setup Script"
echo "====================================="
echo ""

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Please install Go 1.23+"; exit 1; }
command -v git >/dev/null 2>&1 || { echo "❌ Git is not installed."; exit 1; }
command -v dpkg-deb >/dev/null 2>&1 || { echo "❌ dpkg-deb is not installed. Install with: sudo apt install dpkg"; exit 1; }

# Create working directory
WORK_DIR=$(mktemp -d)
echo "📁 Working directory: ${WORK_DIR}"
cd "$WORK_DIR"

# Clone repository
echo ""
echo "📥 Cloning repository..."
git clone "$REPO_URL"
cd "$REPO_NAME"

# Build binaries
echo ""
echo "🔨 Building binaries for Linux..."
mkdir -p build

LDFLAGS="-ldflags \"-s -w -X github.com/Hermithic/aiask/internal/cli.Version=${VERSION}\""

echo "   Building linux-amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X github.com/Hermithic/aiask/internal/cli.Version=${VERSION}" -o build/aiask-linux-amd64 ./cmd/aiask

echo "   Building linux-arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X github.com/Hermithic/aiask/internal/cli.Version=${VERSION}" -o build/aiask-linux-arm64 ./cmd/aiask

# Build .deb packages
echo ""
echo "📦 Building .deb packages..."

for ARCH in amd64 arm64; do
    echo "   Building aiask_${VERSION}_${ARCH}.deb..."
    
    DEB_DIR="build/deb/aiask_${VERSION}_${ARCH}"
    mkdir -p "${DEB_DIR}/DEBIAN"
    mkdir -p "${DEB_DIR}/usr/local/bin"
    
    # Copy binary
    if [ "$ARCH" = "amd64" ]; then
        cp build/aiask-linux-amd64 "${DEB_DIR}/usr/local/bin/aiask"
    else
        cp build/aiask-linux-arm64 "${DEB_DIR}/usr/local/bin/aiask"
    fi
    chmod 755 "${DEB_DIR}/usr/local/bin/aiask"
    
    # Create control file
    cat > "${DEB_DIR}/DEBIAN/control" << EOF
Package: aiask
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: Hermithic <hermithic@users.noreply.github.com>
Description: AI-powered command line assistant
 AIask is an AI-powered command line assistant that converts
 natural language into shell commands for PowerShell, CMD,
 Bash, and Zsh.
 .
 Features:
 - Natural language to command conversion
 - Multiple LLM provider support (Grok, OpenAI, Anthropic, Gemini, Ollama)
 - Auto-detect current shell and OS
 - Execute, copy, edit, or re-prompt commands
 .
 New in v2.0:
 - Command history with search
 - Save and reuse prompt templates
 - Explain mode - understand what commands do
 - Interactive REPL mode
 - Dangerous command warnings and undo suggestions
Homepage: https://github.com/Hermithic/aiask
EOF

    # Create postinst script
    cat > "${DEB_DIR}/DEBIAN/postinst" << 'EOF'
#!/bin/bash
echo ""
echo "✅ AIask installed successfully!"
echo ""
echo "Run 'aiask config' to set up your AI provider."
echo ""
EOF
    chmod 755 "${DEB_DIR}/DEBIAN/postinst"
    
    # Build .deb
    dpkg-deb --build "${DEB_DIR}"
    mv "${DEB_DIR}.deb" "build/aiask_${VERSION}_${ARCH}.deb"
done

# Create APT repository structure
echo ""
echo "📚 Creating APT repository structure..."

APT_REPO="build/apt-repo"
mkdir -p "${APT_REPO}/pool/main/a/aiask"
mkdir -p "${APT_REPO}/dists/stable/main/binary-amd64"
mkdir -p "${APT_REPO}/dists/stable/main/binary-arm64"

# Copy .deb files
cp build/aiask_${VERSION}_amd64.deb "${APT_REPO}/pool/main/a/aiask/"
cp build/aiask_${VERSION}_arm64.deb "${APT_REPO}/pool/main/a/aiask/"

# Generate Packages files
echo "   Generating Packages files..."
cd "${APT_REPO}"

# amd64
dpkg-scanpackages --arch amd64 pool/ > dists/stable/main/binary-amd64/Packages
gzip -k -f dists/stable/main/binary-amd64/Packages

# arm64
dpkg-scanpackages --arch arm64 pool/ > dists/stable/main/binary-arm64/Packages
gzip -k -f dists/stable/main/binary-arm64/Packages

# Generate Release file
echo "   Generating Release file..."
cat > dists/stable/Release << EOF
Origin: Hermithic
Label: AIask
Suite: stable
Codename: stable
Version: ${VERSION}
Architectures: amd64 arm64
Components: main
Description: AIask APT Repository
Date: $(date -Ru)
EOF

# Add checksums to Release
{
    echo "MD5Sum:"
    find dists/stable/main -type f -name "Packages*" -exec sh -c 'echo " $(md5sum {} | cut -d" " -f1) $(stat -c%s {}) {}"' \;
    echo "SHA256:"
    find dists/stable/main -type f -name "Packages*" -exec sh -c 'echo " $(sha256sum {} | cut -d" " -f1) $(stat -c%s {}) {}"' \;
} >> dists/stable/Release

cd "$WORK_DIR/$REPO_NAME"

# Deploy to gh-pages
echo ""
echo "🚀 Deploying to GitHub Pages..."

# Check if gh-pages branch exists
if git ls-remote --heads origin gh-pages | grep -q gh-pages; then
    git fetch origin gh-pages
    git worktree add gh-pages-deploy origin/gh-pages
else
    git worktree add gh-pages-deploy --orphan gh-pages
    cd gh-pages-deploy
    git rm -rf . 2>/dev/null || true
    cd ..
fi

# Copy APT repo to gh-pages
cp -r build/apt-repo/* gh-pages-deploy/

# Also copy .deb files to root for direct download
cp build/aiask_${VERSION}_amd64.deb gh-pages-deploy/
cp build/aiask_${VERSION}_arm64.deb gh-pages-deploy/

# Create index.html
cat > gh-pages-deploy/index.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>AIask APT Repository</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        code { background: #f4f4f4; padding: 2px 6px; border-radius: 4px; }
        pre { background: #1e1e1e; color: #d4d4d4; padding: 16px; border-radius: 8px; overflow-x: auto; }
        h1 { color: #333; }
        a { color: #0366d6; }
    </style>
</head>
<body>
    <h1>🤖 AIask APT Repository</h1>
    <p>AI-powered command line assistant - turn plain English into shell commands!</p>
    
    <h2>Installation</h2>
    <pre>
# Add repository
echo "deb [trusted=yes] https://hermithic.github.io/aiask/ stable main" | sudo tee /etc/apt/sources.list.d/aiask.list

# Install
sudo apt update
sudo apt install aiask</pre>

    <h2>Direct Download</h2>
    <ul>
        <li><a href="aiask_${VERSION}_amd64.deb">aiask_${VERSION}_amd64.deb</a> (x86_64)</li>
        <li><a href="aiask_${VERSION}_arm64.deb">aiask_${VERSION}_arm64.deb</a> (ARM64)</li>
    </ul>

    <h2>Links</h2>
    <ul>
        <li><a href="https://github.com/Hermithic/aiask">GitHub Repository</a></li>
        <li><a href="https://github.com/Hermithic/aiask/releases">All Releases</a></li>
    </ul>
</body>
</html>
EOF

# Commit and push
cd gh-pages-deploy
git add -A
git commit -m "Update APT repository to v${VERSION}"
git push origin gh-pages

cd "$WORK_DIR/$REPO_NAME"
git worktree remove gh-pages-deploy

# Upload .deb to GitHub release
echo ""
echo "📤 Uploading .deb packages to GitHub release..."
gh release upload "v${VERSION}" "build/aiask_${VERSION}_amd64.deb" "build/aiask_${VERSION}_arm64.deb" --clobber

# Cleanup
echo ""
echo "🧹 Cleaning up..."
cd /
rm -rf "$WORK_DIR"

echo ""
echo "✅ APT Repository setup complete!"
echo ""
echo "Users can now install with:"
echo '  echo "deb [trusted=yes] https://hermithic.github.io/aiask/ stable main" | sudo tee /etc/apt/sources.list.d/aiask.list'
echo "  sudo apt update"
echo "  sudo apt install aiask"
echo ""

