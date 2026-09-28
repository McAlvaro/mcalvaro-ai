package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `mcalvaro-ai (Gentle AI Lite) — Ecosystem, Frameworks, Workflows (%s)

USAGE
  mcalvaro-ai                   Launch interactive TUI
  mcalvaro-ai <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine (Clean Triad: Orchestrator, Architect, Builder)
  uninstall    Remove managed files from this machine
  sync         Sync agent configs and skills to current version
  test         Run deterministic test harness (compiler, linter, tests)
  skill-registry refresh
               Refresh .atl/skill-registry.md with cache-hit fast path
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  doctor       Run ecosystem health diagnostics
  version      Print version

FLAGS
  --help, -h    Show global help

Run 'mcalvaro-ai help' for this message.
`, version)
}
