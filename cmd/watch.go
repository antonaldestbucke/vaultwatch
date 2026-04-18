package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

)

var watchInterval.Duration

varCmd = &cobra.Command{
se:   "watch [ort: "Continuously poll Vault and report diffs when they change",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		cfgFile, _ := cmd.Flags().GetString("config")

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		fmt.Fprintf(os.Stdout, "Watching %s every %s — press Ctrl+C to stop\n", path, watchInterval)

		wCfg := audit.WatchConfig{
			Interval: watchInterval,
			Path:     path,
			Clients:  clients,
			OnChange: func(results []audit.CompareResult) {
				fmt.Fprintf(os.Stdout, "[%s] change detected:\n", time.Now().Format(time.RFC3339))
				for _, r := range results {
					audit.PrintTextReport(os.Stdout, r.Diff)
				}
			},
		}

		return audit.Watch(ctx, wCfg)
	},
}

func init() {
	watchCmd.Flags().DurationVarP(&watchInterval, "interval", "i", 30*time.Second, "polling interval")
	watchCmd.Flags().String("config", "vaultwatch.yaml", "config file path")
	rootCmd.AddCommand(watchCmd)
}
