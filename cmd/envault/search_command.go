package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envault/internal/vault"
)

func init() {
	var caseSensitive bool
	var maskValues bool
	var keyOnly bool

	searchCmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for keys or values inside the encrypted vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := vault.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("init vault: %w", err)
			}

			opts := vault.SearchOptions{
				Query:         args[0],
				CaseSensitive: caseSensitive,
				MaskValues:    maskValues,
				KeyOnly:       keyOnly,
			}

			results, err := v.Search(opts)
			if err != nil {
				return fmt.Errorf("search: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(os.Stdout, "No matches found.")
				return nil
			}

			for _, r := range results {
				fmt.Fprintf(os.Stdout, "%s=%s\n", r.Key, r.Value)
			}
			return nil
		},
	}

	searchCmd.Flags().BoolVarP(&caseSensitive, "case-sensitive", "c", false, "Enable case-sensitive matching")
	searchCmd.Flags().BoolVarP(&maskValues, "mask", "m", false, "Mask secret values in output")
	searchCmd.Flags().BoolVarP(&keyOnly, "key-only", "k", false, "Match against keys only")

	rootCmd.AddCommand(searchCmd)
}
