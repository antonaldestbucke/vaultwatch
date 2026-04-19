package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLockCmd_AddMissingFlags(t *testing.T) {
	_, err := executeCommand(rootCmd, "lock", "add", "--file", "/tmp/x.json")
	if err == nil {
		t.Error("expected error when --path and --by are missing")
	}
}

func TestLockCmd_AddAndList(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "locks.json")

	_, err := executeCommand(rootCmd, "lock", "add",
		"--file", file,
		"--path", "secret/prod/db",
		"--by", "alice",
		"--reason", "prod freeze",
	)
	if err != nil {
		t.Fatalf("add: %v", err)
	}

	out, err := executeCommand(rootCmd, "lock", "list", "--file", file)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "secret/prod/db") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestLockCmd_ListEmptyFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "locks.json")

	out, err := executeCommand(rootCmd, "lock", "list", "--file", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No locked paths") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
