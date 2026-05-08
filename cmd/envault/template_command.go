package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envault/internal/vault"
)

func init() {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "template <template-file>",
		Short: "Render a template file using secrets from the vault",
		Long: `Substitute $VAR or ${VAR} placeholders in a template file
using key-value pairs from the decrypted .env file.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := vault.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			v, err := vault.New(cfg, "")
			if err != nil {
				return fmt.Errorf("init vault: %w", err)
			}

			result, err := v.RenderTemplate(args[0])
			if err != nil {
				return fmt.Errorf("render template: %w", err)
			}

			if len(result.Missing) > 0 {
				fmt.Fprintf(os.Stderr, "warning: unresolved variables: %s\n",
					strings.Join(result.Missing, ", "))
			}

			if outputPath == "" || outputPath == "-" {
				fmt.Print(result.Output)
				return nil
			}

			if err := os.WriteFile(outputPath, []byte(result.Output), 0644); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			fmt.Fprintf(os.Stderr, "rendered %d substitution(s) → %s\n", result.Substituted, outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "-", "output file path (default: stdout)")
	rootCmd.AddCommand(cmd)
}
