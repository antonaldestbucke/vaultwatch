package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourorg/vaultwatch/internal/vault"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultwatch.yaml")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestLoadConfig_Valid(t *testing.T) {
	yaml := `
environments:
  - name: dev
    addr: http://vault-dev:8200
    token: dev-token
    mount_path: secret
  - name: prod
    addr: http://vault-prod:8200
    token: prod-token
    mount_path: secret
`
	p := writeConfig(t, yaml)
	cfg, err := vault.LoadConfig(p)
	require.NoError(t, err)
	assert.Len(t, cfg.Environments, 2)
	assert.Equal(t, "dev", cfg.Environments[0].Name)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := vault.LoadConfig("/nonexistent/path.yaml")
	assert.Error(t, err)
}

func TestLoadConfig_NoEnvironments(t *testing.T) {
	p := writeConfig(t, "environments: []\n")
	_, err := vault.LoadConfig(p)
	assert.ErrorContains(t, err, "at least one environment")
}

func TestLoadConfig_DuplicateEnv(t *testing.T) {
	yaml := `
environments:
  - name: dev
    addr: http://vault-dev:8200
  - name: dev
    addr: http://vault-dev2:8200
`
	p := writeConfig(t, yaml)
	_, err := vault.LoadConfig(p)
	assert.ErrorContains(t, err, "duplicate environment")
}

func TestLoadConfig_MissingAddr(t *testing.T) {
	yaml := `
environments:
  - name: dev
    token: tok
`
	p := writeConfig(t, yaml)
	_, err := vault.LoadConfig(p)
	assert.ErrorContains(t, err, "missing 'addr'")
}
