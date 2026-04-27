package audit

import (
	"testing"
	"time"
)

func sampleForecastTrend() []ScoredTrendPoint {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []ScoredTrendPoint{
		{Path: "secret/app", Timestamp: base, Score: 10},
		{Path: "secret/app", Timestamp: base.Add(24 * time.Hour), Score: 20},
		{Path: "secret/app", Timestamp: base.Add(48 * time.Hour), Score: 30},
		{Path: "secret/db", Timestamp: base, Score: 80},
		{Path: "secret/db", Timestamp: base.Add(24 * time.Hour), Score: 70},
		{Path: "secret/db", Timestamp: base.Add(48 * time.Hour), Score: 60},
	}
}

func TestBuildForecast_ReturnsBothPaths(t *testing.T) {
	results := BuildForecast(sampleForecastTrend(), 3, 24*time.Hour)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestBuildForecast_TrendingUp(t *testing.T) {
	results := BuildForecast(sampleForecastTrend(), 2, 24*time.Hour)
	var app ForecastResult
	for _, r := range results {
		if r.Path == "secret/app" {
			app = r
		}
	}
	if app.Trending != "up" {
		t.Errorf("expected trending=up for secret/app, got %q", app.Trending)
	}
}

func TestBuildForecast_TrendingDown(t *testing.T) {
	results := BuildForecast(sampleForecastTrend(), 2, 24*time.Hour)
	var db ForecastResult
	for _, r := range results {
		if r.Path == "secret/db" {
			db = r
		}
	}
	if db.Trending != "down" {
		t.Errorf("expected trending=down for secret/db, got %q", db.Trending)
	}
}

func TestBuildForecast_ForecastLength(t *testing.T) {
	results := BuildForecast(sampleForecastTrend(), 5, 24*time.Hour)
	for _, r := range results {
		if len(r.Forecast) != 5 {
			t.Errorf("path %s: expected 5 forecast points, got %d", r.Path, len(r.Forecast))
		}
	}
}

func TestBuildForecast_ScoreClamped(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	highTrend := []ScoredTrendPoint{
		{Path: "secret/x", Timestamp: base, Score: 90},
		{Path: "secret/x", Timestamp: base.Add(24 * time.Hour), Score: 95},
		{Path: "secret/x", Timestamp: base.Add(48 * time.Hour), Score: 99},
	}
	results := BuildForecast(highTrend, 10, 24*time.Hour)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	for _, fp := range results[0].Forecast {
		if fp.Score > 100 || fp.Score < 0 {
			t.Errorf("score %f out of [0,100] bounds", fp.Score)
		}
	}
}

func TestBuildForecast_Empty(t *testing.T) {
	results := BuildForecast(nil, 3, time.Hour)
	if results != nil {
		t.Errorf("expected nil for empty input, got %v", results)
	}
}

func TestBuildForecast_ZeroSteps(t *testing.T) {
	results := BuildForecast(sampleForecastTrend(), 0, time.Hour)
	if results != nil {
		t.Errorf("expected nil for zero steps, got %v", results)
	}
}
