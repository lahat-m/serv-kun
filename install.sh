#!/usr/bin/env bash
set -euo pipefail

REPO="lahat-m/serv-kun"
BIN_DIR="/usr/local/bin"
BIN_NAME="serv-kun"

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH_SUFFIX="linux-amd64" ;;
  aarch64) ARCH_SUFFIX="linux-arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [[ -z "$LATEST" ]]; then
  echo "Could not determine latest release" >&2
  exit 1
fi

BASE_URL="https://github.com/${REPO}/releases/download/${LATEST}"
BIN_FILE="serv-kun-${ARCH_SUFFIX}"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading serv-kun ${LATEST} (${ARCH_SUFFIX})..."
curl -fsSL "${BASE_URL}/${BIN_FILE}" -o "${TMP_DIR}/${BIN_FILE}"
curl -fsSL "${BASE_URL}/checksums.txt" -o "${TMP_DIR}/checksums.txt"

echo "Verifying checksum..."
if ! (cd "$TMP_DIR" && grep " ${BIN_FILE}\$" checksums.txt | sha256sum -c -); then
  echo "Checksum verification failed — refusing to install." >&2
  exit 1
fi

install -m 0755 "${TMP_DIR}/${BIN_FILE}" "${BIN_DIR}/${BIN_NAME}"

echo "Installed → ${BIN_DIR}/${BIN_NAME}"
echo "Run: sudo serv-kun --help"
