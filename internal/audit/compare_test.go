package audit_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vaultwatch/internal/audit"
	"github.com/vaultwatch/internal/vault"
)

func mockVaultWithKeys(t *testing.T, path string, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		payload := map[string]interface{}{
			"data": map[string]interface{}{"keys": keys},
		}
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestComparePathAcrossEnvs_Diff(t *testing.T) {
	srvA := mockVaultWithKeys(t, "secret/app", []string{"DB_HOST", "DB_PASS"})
	defer srvA.Close()
	srvB := mockVaultWithKeys(t, "secret/app", []string{"DB_HOST", "API_KEY"})
	defer srvB.Close()

	clientA, err := vault.NewClient(srvA.URL, "token-a")
	if err != nil {
		t.Fatalf("clientA: %v", err)
	}
	clientB, err := vault.NewClient(srvB.URL, "token-b")
	if err != nil {
		t.Fatalf("clientB: %v", err)
	}

	result, err := audit.ComparePathAcrossEnvs(context.Background(), "secret/app", "staging", "prod", clientA, clientB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.EnvA != "staging" || result.EnvB != "prod" {
		t.Errorf("unexpected env names: %q %q", result.EnvA, result.EnvB)
	}
	if len(result.DiffResult.OnlyInA) != 1 || result.DiffResult.OnlyInA[0] != "secret/app/DB_PASS" {
		t.Errorf("expected OnlyInA=[secret/app/DB_PASS], got %v", result.DiffResult.OnlyInA)
	}
	if len(result.DiffResult.OnlyInB) != 1 || result.DiffResult.OnlyInB[0] != "secret/app/API_KEY" {
		t.Errorf("expected OnlyInB=[secret/app/API_KEY], got %v", result.DiffResult.OnlyInB)
	}
}

func TestComparePathAcrossEnvs_NoDiff(t *testing.T) {
	keys := []string{"FOO", "BAR"}
	srvA := mockVaultWithKeys(t, "secret/app", keys)
	defer srvA.Close()
	srvB := mockVaultWithKeys(t, "secret/app", keys)
	defer srvB.Close()

	clientA, _ := vault.NewClient(srvA.URL, "token-a")
	clientB, _ := vault.NewClient(srvB.URL, "token-b")

	result, err := audit.ComparePathAcrossEnvs(context.Background(), "secret/app", "dev", "prod", clientA, clientB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.DiffResult.OnlyInA) != 0 || len(result.DiffResult.OnlyInB) != 0 {
		t.Errorf("expected no diff, got %+v", result.DiffResult)
	}
}
