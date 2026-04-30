package audit

import (
	"testing"
)

func sampleSignatureReports() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/app/db", Score: 80, Drifted: false},
		{Path: "secret/app/api", Score: 40, Drifted: true},
		{Path: "secret/infra/tls", Score: 55, Drifted: true},
	}
}

func TestSignReports_ReturnsHex(t *testing.T) {
	reports := sampleSignatureReports()
	sig, err := SignReports(reports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig) != 64 {
		t.Errorf("expected 64-char hex string, got len=%d", len(sig))
	}
}

func TestSignReports_Deterministic(t *testing.T) {
	reports := sampleSignatureReports()
	s1, _ := SignReports(reports)
	s2, _ := SignReports(reports)
	if s1 != s2 {
		t.Errorf("expected deterministic signatures, got %q and %q", s1, s2)
	}
}

func TestSignReports_OrderIndependent(t *testing.T) {
	a := []ScoredReport{
		{Path: "secret/a", Score: 90},
		{Path: "secret/b", Score: 50},
	}
	b := []ScoredReport{
		{Path: "secret/b", Score: 50},
		{Path: "secret/a", Score: 90},
	}
	s1, _ := SignReports(a)
	s2, _ := SignReports(b)
	if s1 != s2 {
		t.Errorf("expected order-independent signature, got %q vs %q", s1, s2)
	}
}

func TestSignReports_Empty(t *testing.T) {
	_, err := SignReports(nil)
	if err == nil {
		t.Error("expected error for empty reports")
	}
}

func TestRecordSignature_AppendsEntry(t *testing.T) {
	store := &SignatureStore{}
	reports := sampleSignatureReports()

	if err := RecordSignature(store, reports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(store.Entries))
	}
	entry := store.Entries[0]
	if entry.PathCount != 3 {
		t.Errorf("expected PathCount=3, got %d", entry.PathCount)
	}
	if entry.DriftCount != 2 {
		t.Errorf("expected DriftCount=2, got %d", entry.DriftCount)
	}
	if entry.Signature == "" {
		t.Error("expected non-empty signature")
	}
}

func TestVerifySignature_Match(t *testing.T) {
	reports := sampleSignatureReports()
	sig, _ := SignReports(reports)

	ok, err := VerifySignature(reports, sig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected signature to match")
	}
}

func TestVerifySignature_Mismatch(t *testing.T) {
	reports := sampleSignatureReports()
	ok, err := VerifySignature(reports, "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected signature mismatch")
	}
}
