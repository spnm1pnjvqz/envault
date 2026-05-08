package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envault/internal/vault"
)

func init() {
	var overwrite bool
	var keys string

	copyCmd := &cobra.Command{
		Use:   "copy <src-env-file> <dst-env-file>",
		Short: "Copy secrets from one env file to another",
		Long: `Copy keys from a source .env file into a destination .env file.

By default existing keys in the destination are preserved. Use --overwrite to
replace them. Optionally restrict which keys are copied with --keys.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			var allowlist []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if t := strings.TrimSpace(k); t != "" {
						allowlist = append(allowlist, t)
					}
				}
			}

			result, err := vault.CopyKeys(src, dst, vault.CopyOptions{
				Overwrite: overwrite,
				Keys:      allowlist,
			})
			if err != nil {
				return fmt.Errorf("copy failed: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Copied:    %d key(s)\n", len(result.Copied))
			fmt.Fprintf(cmd.OutOrStdout(), "Overwrote: %d key(s)\n", len(result.Overwrote))
			fmt.Fprintf(cmd.OutOrStdout(), "Skipped:   %d key(s)\n", len(result.Skipped))

			if len(result.Skipped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Hint: use --overwrite to replace existing keys.\n")
			}
			return nil
		},
	}

	copyCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in the destination")
	copyCmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of keys to copy (default: all)")

	rootCmd.AddCommand(copyCmd)
}
