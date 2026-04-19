package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleLockStore() LockStore {
	return LockStore{
		Locks: []LockEntry{
			{Path: "secret/prod/db", LockedBy: "alice", Reason: "prod freeze", LockedAt: time.Now()},
		},
	}
}

func TestSaveAndLoadLocks_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "locks.json")
	store := sampleLockStore()
	if err := SaveLocks(p, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadLocks(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Locks) != 1 || loaded.Locks[0].Path != "secret/prod/db" {
		t.Errorf("unexpected locks: %+v", loaded.Locks)
	}
}

func TestLoadLocks_MissingFile(t *testing.T) {
	store, err := LoadLocks("/nonexistent/locks.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.Locks) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestLoadLocks_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "locks.json")
	os.WriteFile(p, []byte("not json"), 0644)
	_, err := LoadLocks(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApplyLocks_AnnotatesDriftedLockedPath(t *testing.T) {
	store := sampleLockStore()
	reports := []CompareReport{
		{Path: "secret/prod/db", OnlyInA: []string{"password"}, OnlyInB: []string{}},
		{Path: "secret/dev/db", OnlyInA: []string{"password"}, OnlyInB: []string{}},
	}
	result := ApplyLocks(reports, store)
	if result[0].Notes == "" {
		t.Error("expected locked path to be annotated")
	}
	if result[1].Notes != "" {
		t.Error("expected unlocked path to have no annotation")
	}
}

func TestApplyLocks_NoDriftNoAnnotation(t *testing.T) {
	store := sampleLockStore()
	reports := []CompareReport{
		{Path: "secret/prod/db", OnlyInA: []string{}, OnlyInB: []string{}},
	}
	result := ApplyLocks(reports, store)
	if result[0].Notes != "" {
		t.Errorf("expected no annotation for clean locked path, got: %s", result[0].Notes)
	}
}
