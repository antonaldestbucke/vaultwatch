package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand is a test helper to run a cobra command with args.
func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return buf.String(), err
}

func TestReportCmd_MissingPath(t *testing.T) {
	// reportCmd requires --path flag; omitting it should return an error.
	_, err := executeCommand(rootCmd, "report", "--config", "../configs/vaultwatch.example.yaml")
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
	if !strings.Contains(err.Error(), "path") {
		t.Errorf("expected error to mention 'path', got: %v", err)
	}
}

func TestReportCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "report",
		"--config", "/nonexistent/path.yaml",
		"--path", "secret/app",
	)
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}
