package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "snap.json")

	orig := Snapshot{
		Path:       "secret/myapp",
		Env:        "production",
		Keys:       []string{"db_pass", "api_key"},
		CapturedAt: time.Now().UTC().Truncate(time.Second),
		Meta:       map[string]string{"author": "ci"},
	}

	if err := SaveSnapshot(orig, path); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if loaded.Path != orig.Path {
		t.Errorf("Path: got %q want %q", loaded.Path, orig.Path)
	}
	if loaded.Env != orig.Env {
		t.Errorf("Env: got %q want %q", loaded.Env, orig.Env)
	}
	if len(loaded.Keys) != len(orig.Keys) {
		t.Errorf("Keys length: got %d want %d", len(loaded.Keys), len(orig.Keys))
	}
	if loaded.Meta["author"] != "ci" {
		t.Errorf("Meta author: got %q want %q", loaded.Meta["author"], "ci")
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	s := Snapshot{Path: "secret/x", Env: "dev", Keys: []string{}}
	err := SaveSnapshot(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
