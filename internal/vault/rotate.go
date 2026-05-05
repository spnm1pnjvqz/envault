package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateKeys generates a new key pair, re-encrypts the vault with the new keys,
// and archives the old private key with a timestamp suffix.
func (v *Vault) RotateKeys() error {
	// Load existing config
	cfg, err := LoadConfig(v.configPath)
	if err != nil {
		return fmt.Errorf("rotate: load config: %w", err)
	}

	// Generate new key pair
	recipient, identity, err := GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("rotate: generate key pair: %w", err)
	}

	// Decrypt current vault contents using old keys
	envEntries, err := v.View()
	if err != nil {
		return fmt.Errorf("rotate: decrypt existing vault: %w", err)
	}

	// Archive old private key
	if cfg.PrivateKeyPath != "" {
		expanded := expandHome(cfg.PrivateKeyPath)
		archivePath := expanded + "." + time.Now().Format("20060102150405") + ".bak"
		if err := archiveFile(expanded, archivePath); err != nil {
			return fmt.Errorf("rotate: archive old key: %w", err)
		}
	}

	// Save new key pair
	if err := SaveKeyPair(recipient, identity, cfg.PrivateKeyPath); err != nil {
		return fmt.Errorf("rotate: save new key pair: %w", err)
	}

	// Update config with new public key
	cfg.PublicKey = recipient.String()
	if err := SaveConfig(cfg, v.configPath); err != nil {
		return fmt.Errorf("rotate: save config: %w", err)
	}

	// Re-create vault with new keys
	newVault, err := New(v.configPath)
	if err != nil {
		return fmt.Errorf("rotate: create new vault: %w", err)
	}

	// Re-encrypt with new keys
	if err := newVault.Lock(envEntries); err != nil {
		return fmt.Errorf("rotate: re-encrypt vault: %w", err)
	}

	return nil
}

func archiveFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0600)
}
