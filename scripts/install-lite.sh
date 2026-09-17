#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# Gentle-AI Lite — Fast Install Script
# Clean Triad: Orchestrator + Architect + Builder with Deterministic Test Harness
#
# Usage:
#   curl -sL https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install-lite.sh | bash
# ============================================================================

setup_colors() {
    if [ -t 1 ] && [ "${TERM:-}" != "dumb" ]; then
        RED='\033[0;31m'
        GREEN='\033[0;32m'
        YELLOW='\033[1;33m'
        BLUE='\033[0;34m'
        CYAN='\033[0;36m'
        BOLD='\033[1m'
        NC='\033[0m'
    else
        RED='' GREEN='' YELLOW='' BLUE='' CYAN='' BOLD='' NC=''
    fi
}

setup_colors

info()    { echo -e "${BLUE}[info]${NC}    $*"; }
success() { echo -e "${GREEN}[ok]${NC}      $*"; }
error()   { echo -e "${RED}[error]${NC}   $*" >&2; }
step()    { echo -e "\n${CYAN}${BOLD}==>${NC} ${BOLD}$*${NC}"; }

step "Installing Gentle AI Lite"

INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

if command -v go >/dev/null 2>&1; then
    info "Go compiler detected. Building gentle-ai-lite..."
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

    if [ -f "${REPO_DIR}/go.mod" ]; then
        info "Building from local repository: ${REPO_DIR}"
        go build -o "${INSTALL_DIR}/gentle-ai-lite" "${REPO_DIR}/cmd/gentle-ai"
    else
        info "Installing latest from Go module..."
        GOBIN="${INSTALL_DIR}" go install github.com/gentleman-programming/gentle-ai/v3/cmd/gentle-ai@latest
        mv "${INSTALL_DIR}/gentle-ai" "${INSTALL_DIR}/gentle-ai-lite"
    fi
    ln -sf "${INSTALL_DIR}/gentle-ai-lite" "${INSTALL_DIR}/gentle-ai"
    success "Installed gentle-ai-lite and gentle-ai to ${INSTALL_DIR}"
else
    error "Go compiler not found. Please install Go (1.23+) to build gentle-ai-lite."
    exit 1
fi

step "Verifying Installation"
if "${INSTALL_DIR}/gentle-ai-lite" --version >/dev/null 2>&1; then
    VERSION=$("${INSTALL_DIR}/gentle-ai-lite" --version)
    success "Gentle AI Lite installed successfully: ${VERSION}"
    info "Make sure ${INSTALL_DIR} is in your PATH."
    info "Run 'gentle-ai-lite test' in any project to execute the deterministic test harness."
else
    error "Installation verification failed."
    exit 1
fi
