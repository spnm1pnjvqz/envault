package main

import (
	"fmt"

	"github.com/nicholasgasior/envault/internal/vault"
	"github.com/urfave/cli/v2"
)

var rotateCommand = &cli.Command{
	Name:  "rotate",
	Usage: "Rotate encryption keys and re-encrypt the vault",
	Description: `Generates a new age key pair, re-encrypts the vault file with the
new public key, and archives the old private key with a timestamp suffix.

The old private key is saved as <key-path>.<timestamp>.bak for recovery purposes.
Ensure the vault is locked (encrypted) before running this command.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   ".envault.toml",
			Usage:   "path to envault config file",
			EnvVars: []string{"ENVAULT_CONFIG"},
		},
	},
	Action: func(ctx *cli.Context) error {
		cfgPath := ctx.String("config")

		v, err := vault.New(cfgPath)
		if err != nil {
			return fmt.Errorf("load vault: %w", err)
		}

		fmt.Println("Rotating encryption keys...")

		if err := v.RotateKeys(); err != nil {
			return fmt.Errorf("rotate keys: %w", err)
		}

		fmt.Println("✓ Keys rotated successfully.")
		fmt.Println("  Old private key archived with .bak extension.")
		fmt.Println("  Vault re-encrypted with new public key.")

		return nil
	},
}
