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
echo "  Look News — Instalador"
echo "══════════════════════════════════════"
echo ""

echo "[1/5] Buscando a última versão..."
RELEASE_JSON=$(curl -fsSL "$API_URL")
ZIP_URL=$(echo "$RELEASE_JSON" | grep -o '"browser_download_url": *"[^"]*\.zip"' | head -n 1 | sed -E 's/.*"(https[^"]+)"/\1/')

if [ -z "$ZIP_URL" ]; then
  echo "  ✗ Não encontrei um .zip na última release."
  echo "    Confira em: https://github.com/$REPO/releases"
  exit 1
fi
echo "  ✓ Versão encontrada"

TMP_DIR=$(mktemp -d)
echo "[2/5] Baixando..."
curl --progress-bar -L -o "$TMP_DIR/look-news.zip" "$ZIP_URL"
echo "  ✓ Download concluído"

echo "[3/5] Extraindo..."
(
  unzip -q "$TMP_DIR/look-news.zip" -d "$TMP_DIR"
) &
spinner $! "Extraindo arquivos..."
APP_PATH=$(find "$TMP_DIR" -maxdepth 2 -name "*.app" | head -n 1)

if [ -z "$APP_PATH" ]; then
  echo "  ✗ Não encontrei o .app dentro do zip."
  rm -rf "$TMP_DIR"
  exit 1
fi

echo "[4/5] Instalando em /Applications..."
(
  rm -rf "/Applications/Look News.app"
  cp -R "$APP_PATH" /Applications/
) &
spinner $! "Copiando o aplicativo..."

echo "[5/5] Limpando arquivos temporários..."
rm -rf "$TMP_DIR"
echo "  ✓ Limpeza concluída"

echo ""
echo "══════════════════════════════════════"
echo "  Pronto! Abra 'Look News' pelo"
echo "  Launchpad ou Spotlight."
echo "══════════════════════════════════════"