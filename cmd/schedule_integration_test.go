package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vaultwatch/internal/audit"
)

func TestScheduleCmd_ListMissingFile(t *testing.T) {
	out, err := executeCommand(rootCmd, "schedule", "list", "--file", "/nonexistent.json")
	if err == nil {
		t.Errorf("expected error, got: %s", out)
	}
}

func TestScheduleCmd_AddAndList(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "sched.json")
	_, err := executeCommand(rootCmd, "schedule", "add",
		"--file", f,
		"--name", "test",
		"--path", "secret/test",
		"--interval", "1h",
	)
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	out, err := executeCommand(rootCmd, "schedule", "list", "--file", f)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "test") {
		t.Errorf("expected 'test' in output, got: %s", out)
	}
}

func TestScheduleCmd_AddMissingFlags(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "sched.json")
	_, err := executeCommand(rootCmd, "schedule", "add", "--file", f, "--name", "only-name")
	if err == nil {
		t.Error("expected error for missing flags")
	}
}

func TestScheduleCmd_ListShowsNextDue(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "sched.json")
	store := audit.ScheduleStore{
		Entries: []audit.ScheduleEntry{
			{Name: "e1", Path: "secret/e1", Interval: "2h", Enabled: true},
		},
	}
	data, _ := json.MarshalIndent(store, "", "  ")
	os.WriteFile(f, data, 0644)
	out, err := executeCommand(rootCmd, "schedule", "list", "--file", f)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "e1") {
		t.Errorf("expected entry in output: %s", out)
	}
}
