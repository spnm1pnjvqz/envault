package main

import (
	"fmt"

	"github.com/nicholasgasior/envault/internal/vault"
	"github.com/spf13/cobra"
)

func init() {
	var strategy string

	mergeCmd := &cobra.Command{
		Use:   "merge <dst-env> <src-env>",
		Short: "Merge key/value pairs from src into dst env file",
		Long: `Merge reads key/value pairs from src and writes missing keys into dst.

Conflict strategies:
  ours   – keep the destination value (default)
  theirs – overwrite with the source value
  error  – abort if any conflict is detected`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dstPath := args[0]
			srcPath := args[1]

			var s vault.MergeStrategy
			switch strategy {
			case "ours":
				s = vault.MergeStrategyOurs
			case "theirs":
				s = vault.MergeStrategyTheirs
			case "error":
				s = vault.MergeStrategyError
			default:
				return fmt.Errorf("unknown strategy %q: choose ours, theirs, or error", strategy)
			}

			res, err := vault.MergeEnvFiles(dstPath, srcPath, s)
			if err != nil {
				return err
			}

			if len(res.Added) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Added:       %v\n", res.Added)
			}
			if len(res.Overwritten) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Overwritten: %v\n", res.Overwritten)
			}
			if len(res.Skipped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Skipped:     %v\n", res.Skipped)
			}
			if len(res.Added) == 0 && len(res.Overwritten) == 0 && len(res.Skipped) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to merge.")
			}
			return nil
		},
	}

	mergeCmd.Flags().StringVar(&strategy, "strategy", "ours",
		"Conflict resolution strategy: ours | theirs | error")

	rootCmd.AddCommand(mergeCmd)
}
