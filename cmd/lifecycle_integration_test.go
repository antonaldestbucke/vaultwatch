package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLifecycleCmd_SetMissingPath(t *testing.T) {
	args := []string{"lifecycle", "set", "--stage", "active"}
	out, err := executeCommand(rootCmd, args...)
	if err == nil {
		t.Fatalf("expected error, got output: %s", out)
	}
	if !strings.Contains(err.Error(), "--path") {
		t.Errorf("expected --path error, got: %v", err)
	}
}

func TestLifecycleCmd_SetMissingStage(t *testing.T) {
	args := []string{"lifecycle", "set", "--path", "secret/app/db"}
	_, err := executeCommand(rootCmd, args...)
	if err == nil {
		t.Fatal("expected error for missing --stage")
	}
	if !strings.Contains(err.Error(), "--stage") {
		t.Errorf("expected --stage error, got: %v", err)
	}
}

func TestLifecycleCmd_SetInvalidStage(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "lifecycle.json")
	args := []string{"lifecycle", "set", "--path", "secret/app/db", "--stage", "unknown", "--lifecycle-file", file}
	_, err := executeCommand(rootCmd, args...)
	if err == nil {
		t.Fatal("expected error for invalid stage")
	}
	if !strings.Contains(err.Error(), "invalid stage") {
		t.Errorf("expected invalid stage error, got: %v", err)
	}
}

func TestLifecycleCmd_SetAndList(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "lifecycle.json")
	_, err := executeCommand(rootCmd, "lifecycle", "set",
		"--path", "secret/app/db",
		"--stage", "deprecated",
		"--note", "will be removed",
		"--lifecycle-file", file,
	)
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("lifecycle file not created: %v", err)
	}
	out, err := executeCommand(rootCmd, "lifecycle", "list", "--lifecycle-file", file)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "secret/app/db") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "deprecated") {
		t.Errorf("expected stage in output, got: %s", out)
	}
}
