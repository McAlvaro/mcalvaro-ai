#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# mcalvaro-ai (Gentle AI Lite) — Fast Autonomous Installer (Zero-Brew)
# Clean Triad: Orchestrator + Architect + Builder with Deterministic Test Harness
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

step "Installing mcalvaro-ai (Gentle AI Lite)"

INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

if command -v go >/dev/null 2>&1; then
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

    if [ -f "${REPO_DIR}/go.mod" ]; then
        info "Compiling mcalvaro-ai from repository: ${REPO_DIR}"
        go build -o "${INSTALL_DIR}/mcalvaro-ai" "${REPO_DIR}/cmd/gentle-ai"
    else
        info "Installing from Go module..."
        GOBIN="${INSTALL_DIR}" go install github.com/gentleman-programming/gentle-ai/v3/cmd/gentle-ai@latest
        mv "${INSTALL_DIR}/gentle-ai" "${INSTALL_DIR}/mcalvaro-ai"
    fi

    # Create helpful aliases/symlinks
    ln -sf "${INSTALL_DIR}/mcalvaro-ai" "${INSTALL_DIR}/gentle-ai"
    ln -sf "${INSTALL_DIR}/mcalvaro-ai" "${INSTALL_DIR}/gentle-ai-lite"
    ln -sf "${INSTALL_DIR}/mcalvaro-ai" "${INSTALL_DIR}/gentle"

    success "Installed mcalvaro-ai and symlinks (gentle-ai, gentle) to ${INSTALL_DIR}"
else
    error "Go compiler not found. Please install Go (1.23+) to build mcalvaro-ai."
    exit 1
fi

step "Verifying Installation"
if "${INSTALL_DIR}/mcalvaro-ai" --version >/dev/null 2>&1; then
    VERSION=$("${INSTALL_DIR}/mcalvaro-ai" --version)
    success "Installed successfully: ${VERSION}"
    info "Binary available at: ${INSTALL_DIR}/mcalvaro-ai"
    info "Commands available:"
    info "  • mcalvaro-ai test     - Run deterministic test harness"
    info "  • mcalvaro-ai install  - Deploy Clean Triad to agents (OpenCode, Claude, etc.)"
    info "  • mcalvaro-ai sync     - Sync agents and skills"
else
    error "Installation verification failed."
    exit 1
fi
