package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envault/internal/vault"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Display the vault operation audit log",
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().IntP("last", "n", 0, "Show only the last N entries (0 = all)")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, _ []string) error {
	cfg, err := vault.LoadConfig(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logPath := vault.DefaultAuditLogPath(cfg.Dir)
	log, err := vault.LoadAuditLog(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(cmd.OutOrStdout(), "No audit log found.")
			return nil
		}
		return fmt.Errorf("load audit log: %w", err)
	}

	n, _ := cmd.Flags().GetInt("last")
	events := log.Events
	if n > 0 && n < len(events) {
		events = events[len(events)-n:]
	}

	if len(events) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Audit log is empty.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tVAULT FILE\tSUCCESS\tMESSAGE")
	for _, e := range events {
		status := "ok"
		if !e.Success {
			status = "fail"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Operation,
			e.VaultFile,
			status,
			e.Message,
		)
	}
	return w.Flush()
}
