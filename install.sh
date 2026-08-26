#!/usr/bin/env bash
set -e

REPO="clemilsonazevedo/look-news"
API_URL="https://api.github.com/repos/$REPO/releases/latest"

echo "Look News — instalador"
echo ""
echo "Buscando a última versão..."

RELEASE_JSON=$(curl -fsSL "$API_URL")

ZIP_URL=$(echo "$RELEASE_JSON" | grep -o '"browser_download_url": *"[^"]*\.zip"' | head -n 1 | sed -E 's/.*"(https[^"]+)"/\1/')

if [ -z "$ZIP_URL" ]; then
  echo "Não encontrei um .zip na última release. Confira em:"
  echo "https://github.com/$REPO/releases"
  exit 1
fi

TMP_DIR=$(mktemp -d)

echo "Baixando..."
curl -fsSL -o "$TMP_DIR/look-news.zip" "$ZIP_URL"

echo "Extraindo..."
unzip -q "$TMP_DIR/look-news.zip" -d "$TMP_DIR"

APP_PATH=$(find "$TMP_DIR" -maxdepth 2 -name "*.app" | head -n 1)

if [ -z "$APP_PATH" ]; then
  echo "Não encontrei o .app dentro do zip."
  exit 1
fi

echo "Instalando em /Applications..."
rm -rf "/Applications/Look News.app"
cp -R "$APP_PATH" /Applications/

rm -rf "$TMP_DIR"

echo ""
echo "Pronto! Abra 'Look News' pelo Launchpad ou Spotlight."