package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joeshaw/envault/internal/vault"
)

func init() {
	var editor string

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Open the vault in your $EDITOR for interactive editing",
		Long: `Decrypts the vault to a temporary file, opens it in $EDITOR,
then re-encrypts any changes. The temporary file is removed automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}
			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("open vault: %w", err)
			}
			result, err := v.EditEnvFile(vault.EditOptions{Editor: editor})
			if err != nil {
				return err
			}
			if !result.Modified {
				fmt.Fprintln(os.Stderr, "No changes made.")
				return nil
			}
			if len(result.KeysAdded) > 0 {
				fmt.Fprintf(os.Stdout, "Added:   %s\n", strings.Join(result.KeysAdded, ", "))
			}
			if len(result.KeysChanged) > 0 {
				fmt.Fprintf(os.Stdout, "Changed: %s\n", strings.Join(result.KeysChanged, ", "))
			}
			if len(result.KeysRemoved) > 0 {
				fmt.Fprintf(os.Stdout, "Removed: %s\n", strings.Join(result.KeysRemoved, ", "))
			}
			fmt.Fprintln(os.Stdout, "Vault updated.")
			return nil
		},
	}

	cmd.Flags().StringVar(&editor, "editor", "", "Override the editor binary (default: $EDITOR or vi)")
	rootCmd.AddCommand(cmd)
}
