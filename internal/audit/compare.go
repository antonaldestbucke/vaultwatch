package audit

import (
	"context"
	"fmt"

	"vaultwatch/internal/vault"
)

// CompareReport holds the diff result for a single secret path.
type CompareReport struct {
	Path    string   `json:"path"`
	OnlyInA []string `json:"only_in_a"`
	OnlyInB []string `json:"only_in_b"`
	EnvA    string   `json:"env_a"`
	EnvB    string   `json:"env_b"`
	Notes   string   `json:"notes,omitempty"`
}

// ComparePathAcrossEnvs lists keys at path in two vault clients and diffs them.
func ComparePathAcrossEnvs(ctx context.Context, path string, envA string, clientA *vault.Client, envB string, clientB *vault.Client) (CompareReport, error) {
	keysA, err := clientA.ListSecrets(ctx, path)
	if err != nil {
		return CompareReport{}, fmt.Errorf("list %s in %s: %w", path, envA, err)
	}
	keysB, err := clientB.ListSecrets(ctx, path)
	if err != nil {
		return CompareReport{}, fmt.Errorf("list %s in %s: %w", path, envB, err)
	}
	diff := DiffKeys(path, keysA, keysB)
	return CompareReport{
		Path:    path,
		OnlyInA: diff.OnlyInA,
		OnlyInB: diff.OnlyInB,
		EnvA:    envA,
		EnvB:    envB,
	}, nil
}
