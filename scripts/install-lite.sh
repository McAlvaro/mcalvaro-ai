#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# mcalvaro-ai — Fast Autonomous Installer (Zero-Brew)
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

step "Installing mcalvaro-ai"

INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

# Remove any legacy gentle binaries/symlinks
rm -f "${INSTALL_DIR}/gentle" "${INSTALL_DIR}/gentle-ai" "${INSTALL_DIR}/gentle-ai-lite"

if command -v go >/dev/null 2>&1; then
    REPO_DIR=""
    if [ -n "${BASH_SOURCE[0]:-}" ] && [ -f "${BASH_SOURCE[0]}" ]; then
        SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
    fi

    if [ -n "${REPO_DIR}" ] && [ -f "${REPO_DIR}/go.mod" ]; then
        info "Compiling mcalvaro-ai from repository: ${REPO_DIR}"
        SRC_CMD="${REPO_DIR}/cmd/mcalvaro-ai"
        if [ ! -d "${SRC_CMD}" ]; then SRC_CMD="${REPO_DIR}/cmd/gentle-ai"; fi
        go build -o "${INSTALL_DIR}/mcalvaro-ai" "${SRC_CMD}"
    else
        info "Cloning and building from https://github.com/McAlvaro/mcalvaro-ai.git..."
        TMP_SRC=$(mktemp -d)
        git clone --depth 1 https://github.com/McAlvaro/mcalvaro-ai.git "${TMP_SRC}"
        SRC_CMD="./cmd/mcalvaro-ai"
        if [ ! -d "${TMP_SRC}/cmd/mcalvaro-ai" ]; then SRC_CMD="./cmd/gentle-ai"; fi
        (cd "${TMP_SRC}" && go build -o "${INSTALL_DIR}/mcalvaro-ai" "${SRC_CMD}")
        rm -rf "${TMP_SRC}"
    fi

    success "Installed mcalvaro-ai to ${INSTALL_DIR}/mcalvaro-ai"
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
