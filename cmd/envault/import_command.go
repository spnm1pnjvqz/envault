package main

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envault/internal/vault"
	"github.com/spf13/cobra"
)

var (
	importFormat   string
	importOverwrite bool
)

func init() {
	importCmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import secrets from a .env or JSON file into the vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runImport,
	}

	importCmd.Flags().StringVarP(&importFormat, "format", "f", "",
		"Import format: dotenv or json (auto-detected from extension if omitted)")
	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false,
		"Overwrite existing keys with imported values")

	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	srcPath := args[0]

	cfg, err := vault.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	v, err := vault.New(cfg)
	if err != nil {
		return fmt.Errorf("initialising vault: %w", err)
	}

	var fmt_ vault.ImportFormat
	if importFormat != "" {
		fmt_ = vault.ImportFormat(importFormat)
	} else {
		var detectErr error
		fmt_, detectErr = vault.DetectFormat(srcPath)
		if detectErr != nil {
			return fmt.Errorf("detecting format: %w", detectErr)
		}
	}

	result, err := v.Import(srcPath, fmt_, importOverwrite)
	if err != nil {
		return fmt.Errorf("importing secrets: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Import complete:\n")
	fmt.Fprintf(os.Stdout, "  Added:   %d key(s)\n", len(result.Added))
	fmt.Fprintf(os.Stdout, "  Updated: %d key(s)\n", len(result.Updated))
	fmt.Fprintf(os.Stdout, "  Skipped: %d key(s)\n", len(result.Skipped))

	return nil
}
