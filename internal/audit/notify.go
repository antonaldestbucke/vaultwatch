package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NotifyConfig struct {
	WebhookURL string            `json:"webhook_url"`
	Headers    map[string]string `json:"headers,omitempty"`
	TimeoutSec int               `json:"timeout_sec,omitempty"`
}

type NotifyPayload struct {
	Timestamp string         `json:"timestamp"`
	Alerts    []AlertResult  `json:"alerts"`
	Summary   string         `json:"summary"`
}

func SendWebhook(cfg NotifyConfig, alerts []AlertResult) error {
	if cfg.WebhookURL == "" {
		return fmt.Errorf("webhook_url is required")
	}
	if len(alerts) == 0 {
		return nil
	}

	timeout := cfg.TimeoutSec
	if timeout <= 0 {
		timeout = 10
	}

	payload := NotifyPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Alerts:    alerts,
		Summary:   fmt.Sprintf("%d alert(s) triggered", len(alerts)),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest(http.MethodPost, cfg.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
