#!/usr/bin/env bash
set -e

REPO="clemilsonazevedo/look-news"
API_URL="https://api.github.com/repos/$REPO/releases/latest"

spinner() {
  local pid=$1
  local msg=$2
  local spin='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
  local i=0
  while kill -0 "$pid" 2>/dev/null; do
    i=$(( (i + 1) % ${#spin} ))
    printf "\r  %s %s" "${spin:$i:1}" "$msg"
    sleep 0.1
  done
  printf "\r  ✓ %s\n" "$msg"
}

echo "══════════════════════════════════════"
echo "  Look News — Installer"
echo "══════════════════════════════════════"
echo ""

echo "[1/5] Looking for the latest version..."
RELEASE_JSON=$(curl -fsSL "$API_URL")
ZIP_URL=$(echo "$RELEASE_JSON" | grep -o '"browser_download_url": *"[^"]*\.zip"' | head -n 1 | sed -E 's/.*"(https[^"]+)"/\1/')

if [ -z "$ZIP_URL" ]; then
  echo "  ✗ Could not find a .zip in the latest release."
  echo "    Check: https://github.com/$REPO/releases"
  exit 1
fi
echo "  ✓ Version found"

TMP_DIR=$(mktemp -d)
echo "[2/5] Downloading..."
curl --progress-bar -L -o "$TMP_DIR/look-news.zip" "$ZIP_URL"
echo "  ✓ Download completed"

echo "[3/5] Extracting..."
(
  unzip -q "$TMP_DIR/look-news.zip" -d "$TMP_DIR"
) &
spinner $! "Extracting files..."
APP_PATH=$(find "$TMP_DIR" -maxdepth 2 -name "*.app" | head -n 1)

if [ -z "$APP_PATH" ]; then
  echo "  ✗ Could not find the .app inside the zip."
  rm -rf "$TMP_DIR"
  exit 1
fi

echo "[4/5] Installing to /Applications..."
(
  rm -rf "/Applications/Look News.app"
  cp -R "$APP_PATH" /Applications/
) &
spinner $! "Copying the application..."

echo "[5/5] Cleaning up temporary files..."
rm -rf "$TMP_DIR"
echo "  ✓ Cleanup completed"

echo ""
echo "══════════════════════════════════════"
echo "  Done! Open 'Look News' from"
echo "  Launchpad or Spotlight."
echo "══════════════════════════════════════"