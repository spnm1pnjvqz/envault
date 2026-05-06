package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/envault/internal/vault"
)

func init() {
	var format string
	var output string
	var mask bool

	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export decrypted secrets to a file or stdout",
		Long: `Decrypt the vault and export its contents in dotenv or JSON format.

Examples:
  envault export                        # print dotenv to stdout
  envault export --format json          # print JSON to stdout
  envault export --output secrets.env   # write dotenv to file
  envault export --mask                 # mask secret values`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := vault.LoadConfig("")
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			privKeyPath, err := vault.DefaultPrivateKeyPath()
			if err != nil {
				return fmt.Errorf("resolve key path: %w", err)
			}

			v := vault.New(cfg)
			opts := vault.ExportOptions{
				Format:      vault.ExportFormat(format),
				OutputPath:  output,
				MaskSecrets: mask,
			}

			if err := v.Export(privKeyPath, opts); err != nil {
				return fmt.Errorf("export failed: %w", err)
			}

			if output != "" {
				fmt.Printf("✓ Exported secrets to %s\n", output)
			}
			return nil
		},
	}

	exportCmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv or json")
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	exportCmd.Flags().BoolVar(&mask, "mask", false, "Mask secret values in output")

	rootCmd.AddCommand(exportCmd)
}
