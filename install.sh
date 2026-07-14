#!/usr/bin/env bash
set -euo pipefail

REPO="Aeres-u99/codeatlas"
BINARY="codeatlas"
if ! command -v curl >/dev/null 2>&1; then
    echo "Error: curl is required."
    exit 1
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"

case "$OS" in
    linux|darwin) ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

ASSET_NAME="${BINARY}-${OS}-${ARCH}"
echo "==> Looking for latest release..."

DOWNLOAD_URL=$(
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
    grep "browser_download_url" |
    grep "\".*${ASSET_NAME}\"" |
    sed -E 's/.*"([^"]+)".*/\1/' |
    head -n1
)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Unable to locate release asset:"
    echo "    $ASSET_NAME"
    exit 1
fi
TMP="$(mktemp)"
echo "==> Downloading ${ASSET_NAME}..."
curl -fL "$DOWNLOAD_URL" -o "$TMP"
chmod +x "$TMP"
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"
mv "$TMP" "$INSTALL_DIR/$BINARY"

chmod +x "$INSTALL_DIR/$BINARY"
echo
echo "✔ CodeAtlas installed successfully!"
echo
echo "Installed to:"
echo "    $INSTALL_DIR/$BINARY"

mkdir -p "$HOME/.codeatlas/skills"

if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq is required to install CodeAtlas skills."
    exit 1
fi

mkdir -p "$HOME/.codeatlas/skills"

curl -fsSL "https://api.github.com/repos/Aeres-u99/codeatlas/contents/.codeatlas/skills?ref=master" \
| jq -r '.[].download_url' \
| while read -r url; do
    curl -fsSL "$url" -o "$HOME/.codeatlas/skills/$(basename "$url")"
done

for tool in ".claude" ".codex" ".gemini"; do
    SKILLS_DIR="$HOME/$tool/skills"

    if [ ! -d "$SKILLS_DIR" ]; then
        echo "⚠ $tool/skills not found, skipping."
        continue
    fi

    rm -f "$SKILLS_DIR/codeatlas"
    ln -s "$HOME/.codeatlas/skills" "$SKILLS_DIR/codeatlas"

    echo "✓ Linked CodeAtlas skills to $tool"
done

echo
echo "✔ CodeAtlas Skills installed successfully!"
echo
echo "Installed to:"
echo "    $HOME/.codeatlas/skills"
echo
echo "Linked to:"
[ -d "$HOME/.claude" ] && echo "    ~/.claude/skills"
[ -d "$HOME/.codex" ] && echo "    ~/.codex/skills"
[ -d "$HOME/.gemini" ] && echo "    ~/.gemini/skills"

case ":$PATH:" in
    *":$INSTALL_DIR:"*)
        ;;
    *)
        echo
        echo "$INSTALL_DIR is not on your PATH."
        echo

        SHELL_NAME="$(basename "${SHELL:-}")"

        case "$SHELL_NAME" in
            zsh)
                RC="$HOME/.zshrc"
                ;;
            bash)
                RC="$HOME/.bashrc"
                ;;
            fish)
                RC="$HOME/.config/fish/config.fish"
                ;;
            *)
                RC="your shell configuration file"
                ;;
        esac

        echo "Add the following line to $RC:"
        echo
        echo 'export PATH="$HOME/.local/bin:$PATH"'
        ;;
esac

echo
echo "Verify installation:"
echo "    codeatlas --help"
echo
