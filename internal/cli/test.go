package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gentleman-programming/gentle-ai/v3/internal/harness"
)

// RunTest executes the Gentle AI Lite native deterministic test harness.
func RunTest(args []string, stdout io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve current directory: %w", err)
	}
	if len(args) > 0 && args[0] != "" && args[0] != "--help" && args[0] != "-h" {
		dir = args[0]
	}

	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		_, _ = fmt.Fprintln(stdout, "Usage: gentle-ai test [directory]")
		_, _ = fmt.Fprintln(stdout, "Runs the native deterministic test harness (compiler, linter, tests) for the project.")
		return nil
	}

	res, err := harness.Run(context.Background(), dir, stdout, os.Stderr)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("test harness failed with exit code %d", res.ExitCode)
	}
	return nil
}
