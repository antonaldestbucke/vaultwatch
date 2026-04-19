package cmd

import (
	"testing"
)

func TestReplayRecordCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "replay", "record", "--config", "nonexistent.yaml")
	if err == nil {
		t.Error("expected error when --path is missing")
	}
}

func TestReplayShowCmd_MissingAt(t *testing.T) {
	_, err := executeCommand(rootCmd, "replay", "show", "--replay-file", "nonexistent.json")
	if err == nil {
		t.Error("expected error when --at is missing")
	}
}

func TestReplayShowCmd_InvalidAt(t *testing.T) {
	_, err := executeCommand(rootCmd, "replay", "show", "--replay-file", "nonexistent.json", "--at", "not-a-time")
	if err == nil {
		t.Error("expected error for invalid timestamp")
	}
}

func TestReplayShowCmd_MissingFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "replay", "show", "--replay-file", "/nonexistent/replay.json", "--at", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Error("expected error for missing replay file")
	}
}
