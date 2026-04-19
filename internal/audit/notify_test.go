package audit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sampleAlerts() []AlertResult {
	return []AlertResult{
		{Path: "secret/prod/db", Score: 40, Reason: "score below threshold"},
	}
}

func TestSendWebhook_Success(t *testing.T) {
	var received NotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := NotifyConfig{WebhookURL: ts.URL}
	err := SendWebhook(cfg, sampleAlerts())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(received.Alerts) != 1 {
		t.Errorf("expected 1 alert, got %d", len(received.Alerts))
	}
	if received.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSendWebhook_EmptyAlerts(t *testing.T) {
	cfg := NotifyConfig{WebhookURL: "http://unused"}
	err := SendWebhook(cfg, []AlertResult{})
	if err != nil {
		t.Fatalf("expected no error for empty alerts, got %v", err)
	}
}

func TestSendWebhook_MissingURL(t *testing.T) {
	err := SendWebhook(NotifyConfig{}, sampleAlerts())
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}

func TestSendWebhook_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := NotifyConfig{WebhookURL: ts.URL}
	err := SendWebhook(cfg, sampleAlerts())
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSendWebhook_CustomHeaders(t *testing.T) {
	var gotHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Token")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := NotifyConfig{
		WebhookURL: ts.URL,
		Headers:    map[string]string{"X-Token": "secret123"},
	}
	_ = SendWebhook(cfg, sampleAlerts())
	if gotHeader != "secret123" {
		t.Errorf("expected header value 'secret123', got '%s'", gotHeader)
	}
}
