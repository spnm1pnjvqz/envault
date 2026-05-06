package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/envault/internal/env"
	"github.com/yourusername/envault/internal/vault"
)

func init() {
	var maskSecrets bool

	diffCmd := &cobra.Command{
		Use:   "diff <file-a> <file-b>",
		Short: "Show differences between two .env files",
		Long: `Compare two .env files and display added, removed, and changed keys.

Use --mask to hide secret values in the output.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(args[0], args[1], maskSecrets)
		},
	}

	diffCmd.Flags().BoolVar(&maskSecrets, "mask", false, "Mask secret values in output")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(pathA, pathB string, maskSecrets bool) error {
	entriesA, err := env.ReadFile(pathA)
	if err != nil {
		return fmt.Errorf("reading %s: %w", pathA, err)
	}

	entriesB, err := env.ReadFile(pathB)
	if err != nil {
		return fmt.Errorf("reading %s: %w", pathB, err)
	}

	mapA := env.ToMap(entriesA)
	mapB := env.ToMap(entriesB)

	diffs := vault.Diff(mapA, mapB)
	output := vault.FormatDiff(diffs, maskSecrets)

	fmt.Fprintln(os.Stdout, output)
	return nil
}
