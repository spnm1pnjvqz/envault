package vault

import (
	"fmt"

	"github.com/cqroot/envault/internal/env"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey string
	NewKey string
	Found  bool
}

// RenameKey renames an environment variable key within the encrypted vault.
// If overwrite is false and newKey already exists, an error is returned.
// The encrypted file is re-written after the rename.
func (v *Vault) RenameKey(oldKey, newKey string, overwrite bool) (RenameResult, error) {
	if oldKey == "" {
		return RenameResult{}, fmt.Errorf("old key must not be empty")
	}
	if newKey == "" {
		return RenameResult{}, fmt.Errorf("new key must not be empty")
	}
	if oldKey == newKey {
		return RenameResult{}, fmt.Errorf("old key and new key are identical")
	}

	entries, err := v.decryptedEntries()
	if err != nil {
		return RenameResult{}, fmt.Errorf("decrypt: %w", err)
	}

	oldIdx := -1
	newIdx := -1
	for i, e := range entries {
		switch e.Key {
		case oldKey:
			oldIdx = i
		case newKey:
			newIdx = i
		}
	}

	if oldIdx == -1 {
		return RenameResult{OldKey: oldKey, NewKey: newKey, Found: false}, nil
	}

	if newIdx != -1 && !overwrite {
		return RenameResult{}, fmt.Errorf("key %q already exists; use --overwrite to replace it", newKey)
	}

	// Update the key in place.
	entries[oldIdx].Key = newKey

	// Remove duplicate if overwrite is set and newKey existed.
	if newIdx != -1 {
		entries = append(entries[:newIdx], entries[newIdx+1:]...)
	}

	if err := v.encryptAndSave(entries); err != nil {
		return RenameResult{}, fmt.Errorf("encrypt: %w", err)
	}

	return RenameResult{OldKey: oldKey, NewKey: newKey, Found: true}, nil
}

// decryptedEntries decrypts the vault file and returns parsed env entries.
func (v *Vault) decryptedEntries() ([]env.Entry, error) {
	plaintext, err := v.enc.DecryptFile(v.cfg.EncryptedFile)
	if err != nil {
		return nil, err
	}
	return env.Parse(string(plaintext))
}

// encryptAndSave serializes entries and re-encrypts the vault file.
func (v *Vault) encryptAndSave(entries []env.Entry) error {
	plaintext := env.Serialize(entries)
	return v.enc.EncryptToFile([]byte(plaintext), v.cfg.EncryptedFile)
}
