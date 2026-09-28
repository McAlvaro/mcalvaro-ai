package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Step represents a single command to run within the test harness.
type Step struct {
	Name string
	Cmd  string
	Args []string
}

// StackType identifies the detected project technology stack.
type StackType string

const (
	StackGo     StackType = "go"
	StackNode   StackType = "node"
	StackPython StackType = "python"
	StackRust   StackType = "rust"
	StackCustom StackType = "custom"
	StackNone   StackType = "none"
)

// ProjectStack contains the detected stack type and test steps to execute.
type ProjectStack struct {
	Type  StackType
	Steps []Step
}

// Detect inspects the given directory and determines the tech stack and steps.
func Detect(dir string) (*ProjectStack, error) {
	// 1. Go project (go.mod)
	if fileExists(filepath.Join(dir, "go.mod")) {
		steps := []Step{
			{Name: "Go Vet", Cmd: "go", Args: []string{"vet", "./..."}},
			{Name: "Go Test", Cmd: "go", Args: []string{"test", "./..."}},
		}
		return &ProjectStack{Type: StackGo, Steps: steps}, nil
	}

	// 2. Node / TypeScript / JavaScript project (package.json)
	pkgJSONPath := filepath.Join(dir, "package.json")
	if fileExists(pkgJSONPath) {
		pm := detectNodePackageManager(dir)
		steps := detectNodeSteps(dir, pkgJSONPath, pm)
		if len(steps) > 0 {
			return &ProjectStack{Type: StackNode, Steps: steps}, nil
		}
	}

	// 3. Rust project (Cargo.toml)
	if fileExists(filepath.Join(dir, "Cargo.toml")) {
		steps := []Step{
			{Name: "Cargo Check", Cmd: "cargo", Args: []string{"check"}},
			{Name: "Cargo Test", Cmd: "cargo", Args: []string{"test"}},
		}
		return &ProjectStack{Type: StackRust, Steps: steps}, nil
	}

	// 4. Python project (pyproject.toml, pytest.ini, requirements.txt, setup.py)
	if fileExists(filepath.Join(dir, "pyproject.toml")) ||
		fileExists(filepath.Join(dir, "pytest.ini")) ||
		fileExists(filepath.Join(dir, "setup.py")) ||
		fileExists(filepath.Join(dir, "requirements.txt")) {
		steps := []Step{
			{Name: "Python Test", Cmd: "pytest", Args: []string{}},
		}
		return &ProjectStack{Type: StackPython, Steps: steps}, nil
	}

	return &ProjectStack{Type: StackNone, Steps: nil}, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func detectNodePackageManager(dir string) string {
	if fileExists(filepath.Join(dir, "pnpm-lock.yaml")) {
		return "pnpm"
	}
	if fileExists(filepath.Join(dir, "bun.lockb")) || fileExists(filepath.Join(dir, "bun.lock")) {
		return "bun"
	}
	if fileExists(filepath.Join(dir, "yarn.lock")) {
		return "yarn"
	}
	return "npm"
}

func detectNodeSteps(dir, pkgJSONPath, pm string) []Step {
	var steps []Step

	data, err := os.ReadFile(pkgJSONPath)
	if err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			if _, ok := pkg.Scripts["typecheck"]; ok {
				steps = append(steps, Step{Name: "Node Typecheck", Cmd: pm, Args: []string{"run", "typecheck"}})
			} else if fileExists(filepath.Join(dir, "tsconfig.json")) {
				steps = append(steps, Step{Name: "TypeScript Check", Cmd: "npx", Args: []string{"tsc", "--noEmit"}})
			}

			if _, ok := pkg.Scripts["lint"]; ok {
				steps = append(steps, Step{Name: "Node Lint", Cmd: pm, Args: []string{"run", "lint"}})
			}

			if _, ok := pkg.Scripts["test"]; ok {
				if pm == "npm" {
					steps = append(steps, Step{Name: "Node Test", Cmd: "npm", Args: []string{"test", "--", "--run"}})
				} else {
					steps = append(steps, Step{Name: "Node Test", Cmd: pm, Args: []string{"test"}})
				}
			}
		}
	}

	return steps
}

// StepResult holds detailed execution metrics for an individual step.
type StepResult struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	Passed     bool   `json:"passed"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Output     string `json:"output,omitempty"`
}

// Result carries summary information about the executed harness run.
type Result struct {
	Stack       StackType    `json:"stack"`
	Passed      bool         `json:"passed"`
	TotalSteps  int          `json:"total_steps"`
	PassedSteps int          `json:"passed_steps"`
	FailedStep  string       `json:"failed_step,omitempty"`
	ExitCode    int          `json:"exit_code"`
	DurationMs  int64        `json:"duration_ms"`
	Steps       []StepResult `json:"steps"`
}

// Options configures the test harness execution.
type Options struct {
	JSON    bool
	Timeout time.Duration
}

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiDim    = "\033[2m"
)

// Run executes the test harness for the project located at dir with default options.
func Run(ctx context.Context, dir string, stdout, stderr io.Writer) (*Result, error) {
	return RunWithOptions(ctx, dir, stdout, stderr, Options{})
}

// RunWithOptions executes the test harness with the specified configuration options.
func RunWithOptions(ctx context.Context, dir string, stdout, stderr io.Writer, opts Options) (*Result, error) {
	startTime := time.Now()

	stack, err := Detect(dir)
	if err != nil {
		return nil, fmt.Errorf("detect stack: %w", err)
	}

	if opts.Timeout == 0 {
		opts.Timeout = 3 * time.Minute
	}

	res := &Result{
		Stack:      stack.Type,
		Passed:     true,
		TotalSteps: len(stack.Steps),
		Steps:      make([]StepResult, 0, len(stack.Steps)),
	}

	if stack.Type == StackNone || len(stack.Steps) == 0 {
		if opts.JSON {
			res.DurationMs = time.Since(startTime).Milliseconds()
			_ = json.NewEncoder(stdout).Encode(res)
			return res, nil
		}
		_, _ = fmt.Fprintf(stdout, "%sNo supported project stack detected in %s (no go.mod, package.json, Cargo.toml, or Python configs found).%s\n", ansiDim, dir, ansiReset)
		return res, nil
	}

	if !opts.JSON {
		_, _ = fmt.Fprintf(stdout, "\n%s%s==> mcalvaro-ai Test Harness: detected %s stack with %d check(s)%s\n\n", ansiBold, ansiCyan, strings.ToUpper(string(stack.Type)), len(stack.Steps), ansiReset)
	}

	for i, step := range stack.Steps {
		stepStart := time.Now()
		cmdStr := step.Cmd + " " + strings.Join(step.Args, " ")

		if !opts.JSON {
			_, _ = fmt.Fprintf(stdout, "  %s[%d/%d]%s %s%-20s%s %s(%s)%s\n", ansiDim, i+1, len(stack.Steps), ansiReset, ansiBold, step.Name, ansiReset, ansiDim, cmdStr, ansiReset)
		}

		stepCtx, stepCancel := context.WithTimeout(ctx, opts.Timeout)
		cmd := exec.CommandContext(stepCtx, step.Cmd, step.Args...)
		cmd.Dir = dir

		var buf bytes.Buffer
		if opts.JSON {
			cmd.Stdout = &buf
			cmd.Stderr = &buf
		} else {
			// Tee output to buffer and stdout/stderr
			cmd.Stdout = io.MultiWriter(stdout, &buf)
			cmd.Stderr = io.MultiWriter(stderr, &buf)
		}

		cmdErr := cmd.Run()
		stepCancel()
		stepDuration := time.Since(stepStart).Milliseconds()

		stepRes := StepResult{
			Name:       step.Name,
			Command:    cmdStr,
			DurationMs: stepDuration,
			Output:     buf.String(),
		}

		if cmdErr != nil {
			exitCode := 1
			if exitErr, ok := cmdErr.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
			stepRes.Passed = false
			stepRes.ExitCode = exitCode
			res.Steps = append(res.Steps, stepRes)

			res.Passed = false
			res.FailedStep = step.Name
			res.ExitCode = exitCode
			res.DurationMs = time.Since(startTime).Milliseconds()

			if opts.JSON {
				_ = json.NewEncoder(stdout).Encode(res)
				return res, fmt.Errorf("step %q failed: %w", step.Name, cmdErr)
			}

			_, _ = fmt.Fprintf(stderr, "\n  %s✖ FAILED:%s %s (exit code %d) %s[%dms]%s\n", ansiBold+ansiRed, ansiReset, step.Name, exitCode, ansiDim, stepDuration, ansiReset)
			return res, fmt.Errorf("step %q failed: %w", step.Name, cmdErr)
		}

		stepRes.Passed = true
		stepRes.ExitCode = 0
		res.Steps = append(res.Steps, stepRes)
		res.PassedSteps++

		if !opts.JSON {
			_, _ = fmt.Fprintf(stdout, "  %s✔ PASSED:%s %s %s[%dms]%s\n\n", ansiBold+ansiGreen, ansiReset, step.Name, ansiDim, stepDuration, ansiReset)
		}
	}

	res.DurationMs = time.Since(startTime).Milliseconds()

	if opts.JSON {
		_ = json.NewEncoder(stdout).Encode(res)
		return res, nil
	}

	_, _ = fmt.Fprintf(stdout, "%s%s✔ ALL %d CHECKS PASSED (exit code 0)%s %s[%dms]%s\n\n", ansiBold+ansiGreen, ansiBold, len(stack.Steps), ansiReset, ansiDim, res.DurationMs, ansiReset)
	return res, nil
}
