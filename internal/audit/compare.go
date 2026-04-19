package audit

import (
	"context"
	"fmt"

	"vaultwatch/internal/vault"
)

// CompareResult holds the diff result for a single secret path across environments.
type CompareResult struct {
	Path    string   `json:"path"`
	OnlyInA []string `json:"only_in_a"`
	OnlyInB []string `json:"only_in_b"`
	Note    string   `json:"note,omitempty"`
}

// ComparePathAcrossEnvs lists secrets at path in two Vault clients and diffs the keys.
func ComparePathAcrossEnvs(ctx context.Context, path string, clientA, clientB *vault.Client) (CompareResult, error) {
	keysA, err := clientA.ListSecrets(ctx, path)
	if err != nil {
		return CompareResult{}, fmt.Errorf("list secrets in env A at %q: %w", path, err)
	}
	keysB, err := clientB.ListSecrets(ctx, path)
	if err != nil {
		return CompareResult{}, fmt.Errorf("list secrets in env B at %q: %w", path, err)
	}
	diff := DiffKeys(path, keysA, keysB)
	return CompareResult{
		Path:    diff.Path,
		OnlyInA: diff.OnlyInA,
		OnlyInB: diff.OnlyInB,
	}, nil
}
