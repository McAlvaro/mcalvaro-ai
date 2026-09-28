package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gentleman-programming/gentle-ai/v3/internal/harness"
)

// RunTest executes the mcalvaro-ai / Gentle AI Lite native deterministic test harness.
func RunTest(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.SetOutput(stdout)

	jsonOutput := fs.Bool("json", false, "output results in machine-readable JSON format")
	timeout := fs.Duration("timeout", 3*time.Minute, "timeout per verification step (e.g. 2m, 30s)")

	fs.Usage = func() {
		_, _ = fmt.Fprintln(stdout, "Usage: mcalvaro-ai test [flags] [directory]")
		_, _ = fmt.Fprintln(stdout, "\nRuns the native deterministic test harness (compiler, linter, tests) for the project.")
		_, _ = fmt.Fprintln(stdout, "\nFlags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve current directory: %w", err)
	}

	if fs.NArg() > 0 && fs.Arg(0) != "" {
		dir = fs.Arg(0)
	}

	opts := harness.Options{
		JSON:    *jsonOutput,
		Timeout: *timeout,
	}

	res, err := harness.RunWithOptions(context.Background(), dir, stdout, os.Stderr, opts)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("test harness failed with exit code %d", res.ExitCode)
	}
	return nil
}
