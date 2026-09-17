package harness

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDetect_Go(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	stack, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect error: %v", err)
	}
	if stack.Type != StackGo {
		t.Errorf("got stack %v, want %v", stack.Type, StackGo)
	}
	if len(stack.Steps) != 2 {
		t.Errorf("got %d steps, want 2", len(stack.Steps))
	}
}

func TestDetect_Node(t *testing.T) {
	dir := t.TempDir()
	pkgJSON := `{"scripts": {"test": "vitest run", "lint": "eslint ."}}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "pnpm-lock.yaml"), []byte("lockfileVersion: 5.4"), 0644); err != nil {
		t.Fatal(err)
	}

	stack, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect error: %v", err)
	}
	if stack.Type != StackNode {
		t.Errorf("got stack %v, want %v", stack.Type, StackNode)
	}
	if len(stack.Steps) < 2 {
		t.Errorf("expected at least 2 steps (lint, test), got %d", len(stack.Steps))
	}
}

func TestDetect_None(t *testing.T) {
	dir := t.TempDir()
	stack, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect error: %v", err)
	}
	if stack.Type != StackNone {
		t.Errorf("got %v, want %v", stack.Type, StackNone)
	}
}

func TestRun_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	res, err := Run(context.Background(), dir, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
}
