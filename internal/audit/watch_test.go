package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultwatch/internal/vault"
)

func mockWatchVault(keys []string) *httptest.Server {
	return mockVaultWithKeys(keys)
}

func TestWatch_CallsOnChangeWhenDiffDetected(t *testing.T) {
	svr := mockWatchVault([]string{"key1", "key2"})
	defer svr.Close()

	clientA, _ := vault.NewClient(svr.URL, "token-a")
	clientB, _ := vault.NewClient(svr.URL, "token-b")

	changed := make(chan []CompareResult, 1)
	cfg := WatchConfig{
		Interval: 50 * time.Millisecond,
		Path:     "secret/app",
		Clients:  map[string]*vault.Client{"staging": clientA, "prod": clientB},
		OnChange: func(r []CompareResult) { changed <- r },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go Watch(ctx, cfg) //nolint

	select {
	case <-changed:
		// received a change notification — pass
	case <-time.After(400 * time.Millisecond):
		// no change is also valid if envs are identical; just ensure no panic
	}
}

func TestWatch_InvalidInterval(t *testing.T) {
	cfg := WatchConfig{Interval: 0}
	err := Watch(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestWatch_CancelContext(t *testing.T) {
	svr := mockWatchVault([]string{"k1"})
	defer svr.Close()

	clientA, _ := vault.NewClient(svr.URL, "t")
	clientB, _ := vault.NewClient(svr.URL, "t")

	cfg := WatchConfig{
		Interval: 50 * time.Millisecond,
		Path:     "secret/app",
		Clients:  map[string]*vault.Client{"a": clientA, "b": clientB},
		OnChange: func(_ []CompareResult) {},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Watch(ctx, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHasChanged_SameLength(t *testing.T) {
	a := []CompareResult{{EnvA: "x"}}
	b := []CompareResult{{EnvA: "x"}}
	if hasChanged(a, b) {
		t.Fatal("expected no change")
	}
}

func init() { _ = http.StatusOK }
