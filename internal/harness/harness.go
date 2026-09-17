package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		return &ProjectStack{Type: StackNode, Steps: steps}, nil
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

	if len(steps) == 0 {
		steps = append(steps, Step{Name: "Node Test", Cmd: pm, Args: []string{"test"}})
	}

	return steps
}

// Result carries summary information about the executed harness run.
type Result struct {
	Stack       StackType
	TotalSteps  int
	PassedSteps int
	FailedStep  string
	ExitCode    int
}

// Run executes the test harness for the project located at dir.
func Run(ctx context.Context, dir string, stdout, stderr io.Writer) (*Result, error) {
	stack, err := Detect(dir)
	if err != nil {
		return nil, fmt.Errorf("detect stack: %w", err)
	}

	if stack.Type == StackNone || len(stack.Steps) == 0 {
		_, _ = fmt.Fprintf(stdout, "No supported project stack detected in %s (no go.mod, package.json, Cargo.toml, or Python configs found).\n", dir)
		return &Result{Stack: StackNone, ExitCode: 0}, nil
	}

	_, _ = fmt.Fprintf(stdout, "==> Gentle AI Test Harness: detected %s stack with %d step(s)\n", strings.ToUpper(string(stack.Type)), len(stack.Steps))

	res := &Result{
		Stack:      stack.Type,
		TotalSteps: len(stack.Steps),
	}

	for i, step := range stack.Steps {
		_, _ = fmt.Fprintf(stdout, "\n[%d/%d] Running %s (%s %s)...\n", i+1, len(stack.Steps), step.Name, step.Cmd, strings.Join(step.Args, " "))

		cmd := exec.CommandContext(ctx, step.Cmd, step.Args...)
		cmd.Dir = dir
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		if err := cmd.Run(); err != nil {
			exitCode := 1
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
			res.FailedStep = step.Name
			res.ExitCode = exitCode
			_, _ = fmt.Fprintf(stderr, "\nFAILED: step %q failed with exit code %d\n", step.Name, exitCode)
			return res, fmt.Errorf("step %q failed: %w", step.Name, err)
		}

		res.PassedSteps++
		_, _ = fmt.Fprintf(stdout, "PASSED: %s\n", step.Name)
	}

	_, _ = fmt.Fprintf(stdout, "\n==> ALL CHECKS PASSED (exit code 0)\n")
	res.ExitCode = 0
	return res, nil
}
