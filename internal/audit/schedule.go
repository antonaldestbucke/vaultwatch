package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ScheduleEntry struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Interval string `json:"interval"`
	Enabled  bool   `json:"enabled"`
	LastRun  string `json:"last_run,omitempty"`
}

type ScheduleStore struct {
	Entries []ScheduleEntry `json:"entries"`
}

func SaveSchedule(path string, store ScheduleStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schedule: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadSchedule(path string) (ScheduleStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ScheduleStore{}, fmt.Errorf("read schedule: %w", err)
	}
	var store ScheduleStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ScheduleStore{}, fmt.Errorf("parse schedule: %w", err)
	}
	return store, nil
}

func NextDue(entry ScheduleEntry) (time.Duration, error) {
	interval, err := time.ParseDuration(entry.Interval)
	if err != nil {
		return 0, fmt.Errorf("invalid interval %q: %w", entry.Interval, err)
	}
	if entry.LastRun == "" {
		return 0, nil
	}
	last, err := time.Parse(time.RFC3339, entry.LastRun)
	if err != nil {
		return 0, fmt.Errorf("invalid last_run: %w", err)
	}
	next := last.Add(interval)
	return time.Until(next), nil
}
