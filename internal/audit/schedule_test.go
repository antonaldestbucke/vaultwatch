package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleStore() ScheduleStore {
	return ScheduleStore{
		Entries: []ScheduleEntry{
			{Name: "prod-check", Path: "secret/prod", Interval: "1h", Enabled: true},
			{Name: "staging-check", Path: "secret/staging", Interval: "30m", Enabled: false},
		},
	}
}

func TestSaveAndLoadSchedule_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "schedule.json")
	store := sampleStore()
	if err := SaveSchedule(p, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSchedule(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Name != "prod-check" {
		t.Errorf("unexpected name: %s", loaded.Entries[0].Name)
	}
}

func TestLoadSchedule_MissingFile(t *testing.T) {
	_, err := LoadSchedule("/nonexistent/schedule.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSchedule_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	os.WriteFile(p, []byte("not json"), 0644)
	_, err := LoadSchedule(p)
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestNextDue_NoLastRun(t *testing.T) {
	entry := ScheduleEntry{Interval: "1h"}
	d, err := NextDue(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 0 {
		t.Errorf("expected 0 duration, got %v", d)
	}
}

func TestNextDue_WithLastRun(t *testing.T) {
	last := time.Now().Add(-30 * time.Minute).Format(time.RFC3339)
	entry := ScheduleEntry{Interval: "1h", LastRun: last}
	d, err := NextDue(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d <= 0 || d > 31*time.Minute {
		t.Errorf("unexpected duration: %v", d)
	}
}

func TestNextDue_InvalidInterval(t *testing.T) {
	entry := ScheduleEntry{Interval: "bad"}
	_, err := NextDue(entry)
	if err == nil {
		t.Error("expected error for invalid interval")
	}
}
