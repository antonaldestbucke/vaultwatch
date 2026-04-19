package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPolicyCmd_MissingPath(t *testing.T) {
	out, err := executeCommand(rootCmd, "policy", "--policy", "p.json")
	if err == nil {
		t.Errorf("expected error, got output: %s", out)
	}
}

func TestPolicyCmd_MissingPolicy(t *testing.T) {
	out, err := executeCommand(rootCmd, "policy", "--path", "secret/prod")
	if err == nil {
		t.Errorf("expected error, got output: %s", out)
	}
}

func TestPolicyCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "policy",
		"--path", "secret/prod",
		"--policy", "p.json",
		"--config", "/nonexistent/config.yaml",
	)
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestPolicyCmd_InvalidPolicyFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "cfg.yaml")
	_ = os.WriteFile(cfgPath, []byte("environments:\n  - name: prod\n    address: http://localhost:8200\n    token: root\n"), 0644)

	_, err := executeCommand(rootCmd, "policy",
		"--path", "secret/prod",
		"--policy", "/nonexistent/policy.json",
		"--config", cfgPath,
	)
	if err == nil {
		t.Error("expected error for missing policy file")
	}
}
