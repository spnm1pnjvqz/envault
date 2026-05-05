package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := vault.DefaultConfig()
		v, err := vault.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to create vault: %w", err)
		}
		if err := vault.SaveConfig(cfg, ".envault.toml"); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Vault initialized. Public key: %s\n", v.PublicKey())
		fmt.Println("Add .envault.toml to version control. Keep your private key safe!")
		return nil
	},
}

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Encrypt .env file into .env.age",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(".envault.toml")
		if err != nil {
			return fmt.Errorf("failed to load config (run 'envault init' first): %w", err)
		}
		v, err := vault.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to open vault: %w", err)
		}
		if err := v.Lock(); err != nil {
			return fmt.Errorf("lock failed: %w", err)
		}
		fmt.Printf("Encrypted %s -> %s\n", cfg.EnvFile, cfg.EncryptedFile)
		return nil
	},
}

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Decrypt .env.age into .env",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(".envault.toml")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		v, err := vault.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to open vault: %w", err)
		}
		if err := v.Unlock(); err != nil {
			return fmt.Errorf("unlock failed: %w", err)
		}
		fmt.Printf("Decrypted %s -> %s\n", cfg.EncryptedFile, cfg.EnvFile)
		return nil
	},
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View decrypted secrets without writing to disk",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(".envault.toml")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		v, err := vault.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to open vault: %w", err)
		}
		entries, err := v.View()
		if err != nil {
			return fmt.Errorf("view failed: %w", err)
		}
		for _, e := range entries {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
		return nil
	},
}
