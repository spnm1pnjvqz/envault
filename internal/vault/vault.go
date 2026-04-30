// Package vault provides high-level operations for managing encrypted .env vaults.
package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envault/internal/crypto"
	"github.com/user/envault/internal/env"
)

const defaultVaultExt = ".vault"

// Vault manages encryption and decryption of .env files.
type Vault struct {
	encryptor *crypto.Encryptor
}

// New creates a new Vault using the provided age public and private keys.
func New(publicKey, privateKey string) (*Vault, error) {
	enc, err := crypto.NewEncryptor(publicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to initialise encryptor: %w", err)
	}
	return &Vault{encryptor: enc}, nil
}

// Lock encrypts the given .env file and writes the ciphertext to a .vault file.
// The output path is derived from the input path unless outPath is non-empty.
func (v *Vault) Lock(envPath, outPath string) error {
	entries, err := env.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("vault: lock: %w", err)
	}

	plaintext := []byte(env.Serialize(entries))

	if outPath == "" {
		ext := filepath.Ext(envPath)
		outPath = envPath[:len(envPath)-len(ext)] + defaultVaultExt
	}

	if err := v.encryptor.EncryptFile(plaintext, outPath); err != nil {
		return fmt.Errorf("vault: lock: %w", err)
	}
	return nil
}

// Unlock decrypts the given .vault file and writes the plaintext to an .env file.
// The output path is derived from the input path unless outPath is non-empty.
func (v *Vault) Unlock(vaultPath, outPath string) error {
	plaintext, err := v.encryptor.DecryptFile(vaultPath)
	if err != nil {
		return fmt.Errorf("vault: unlock: %w", err)
	}

	if outPath == "" {
		ext := filepath.Ext(vaultPath)
		outPath = vaultPath[:len(vaultPath)-len(ext)] + ".env"
	}

	if err := os.WriteFile(outPath, plaintext, 0o600); err != nil {
		return fmt.Errorf("vault: unlock: write output: %w", err)
	}
	return nil
}

// View decrypts the vault file and returns the parsed env entries without
// writing anything to disk.
func (v *Vault) View(vaultPath string) ([]env.Entry, error) {
	plaintext, err := v.encryptor.DecryptFile(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("vault: view: %w", err)
	}

	entries, err := env.Parse(string(plaintext))
	if err != nil {
		return nil, fmt.Errorf("vault: view: %w", err)
	}
	return entries, nil
}
