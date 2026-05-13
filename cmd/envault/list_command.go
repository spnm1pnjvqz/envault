package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envault/internal/vault"
)

func init() {
	var prefix string
	var maskValues bool
	var noSort bool
	var showTags bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys in the decrypted env file",
		Long: `List displays all key=value pairs from the decrypted .env file.

Use --prefix to filter keys by a common prefix, --mask to hide values,
and --tags to show tag annotations inline.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := vault.LoadConfig(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("init vault: %w", err)
			}

			sortBy := "key"
			if noSort {
				sortBy = "none"
			}

			results, err := v.ListKeys(vault.ListOptions{
				FilterPrefix: prefix,
				MaskValues:   maskValues,
				SortBy:       sortBy,
				ShowTags:     showTags,
			})
			if err != nil {
				return fmt.Errorf("list keys: %w", err)
			}

			fmt.Fprint(os.Stdout, vault.FormatList(results, maskValues))
			return nil
		},
	}

	listCmd.Flags().StringVar(&prefix, "prefix", "", "Filter keys by prefix")
	listCmd.Flags().BoolVar(&maskValues, "mask", false, "Mask secret values")
	listCmd.Flags().BoolVar(&noSort, "no-sort", false, "Preserve original key order")
	listCmd.Flags().BoolVar(&showTags, "tags", false, "Show tag annotations")

	rootCmd.AddCommand(listCmd)
}
