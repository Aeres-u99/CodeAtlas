#!/usr/bin/env bash
set -euo pipefail

REPO="Aeres-u99/hermes"
BINARY="hermes"
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
echo "✔ Hermes installed successfully!"
echo
echo "Installed to:"
echo "    $INSTALL_DIR/$BINARY"

mkdir -p "$HOME/.hermes/skills"

if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq is required to install Hermes skills."
    exit 1
fi

curl -fsSL "https://api.github.com/repos/Aeres-u99/hermes/contents/.hermes/skills?ref=master" \
| jq -r '.[].download_url' \
| while read -r url; do
    curl -fsSL "$url" -o "$HOME/.hermes/skills/$(basename "$url")"
done

for tool in ".claude" ".codex" ".gemini"; do
    if [ -d "$HOME/$tool" ]; then
        if [ ! -e "$HOME/$tool/skills" ]; then
            ln -s "$HOME/.hermes/skills" "$HOME/$tool/skills"
            echo "✓ Linked skills for $tool"
        else
            echo "⚠ $tool/skills already exists, skipping."
        fi
    fi
done

echo
echo "✔ Hermes Skills installed successfully!"
echo
echo "Installed to:"
echo "    $HOME/.hermes/skills"
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
echo "    hermes --help"
echo
