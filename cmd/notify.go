package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send alert results to a webhook endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		alertsFile, _ := cmd.Flags().GetString("alerts")
		webhookURL, _ := cmd.Flags().GetString("webhook")
		token, _ := cmd.Flags().GetString("token")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if alertsFile == "" {
			return fmt.Errorf("--alerts is required")
		}
		if webhookURL == "" {
			return fmt.Errorf("--webhook is required")
		}

		data, err := os.ReadFile(alertsFile)
		if err != nil {
			return fmt.Errorf("failed to read alerts file: %w", err)
		}

		var alerts []audit.AlertResult
		if err := json.Unmarshal(data, &alerts); err != nil {
			return fmt.Errorf("failed to parse alerts file: %w", err)
		}

		cfg := audit.NotifyConfig{
			WebhookURL: webhookURL,
			TimeoutSec: timeout,
		}
		if token != "" {
			cfg.Headers = map[string]string{"Authorization": "Bearer " + token}
		}

		if err := audit.SendWebhook(cfg, alerts); err != nil {
			return fmt.Errorf("webhook delivery failed: %w", err)
		}

		fmt.Printf("Delivered %d alert(s) to %s\n", len(alerts), webhookURL)
		return nil
	},
}

func init() {
	notifyCmd.Flags().String("alerts", "", "Path to JSON file containing alert results")
	notifyCmd.Flags().String("webhook", "", "Webhook URL to POST alerts to")
	notifyCmd.Flags().String("token", "", "Optional bearer token for Authorization header")
	notifyCmd.Flags().Int("timeout", 10, "HTTP timeout in seconds")
	rootCmd.AddCommand(notifyCmd)
}
