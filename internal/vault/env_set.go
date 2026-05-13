package vault

import (
	"fmt"

	"github.com/nicholasgasior/envault/internal/env"
)

// SetResult holds the outcome of a Set operation.
type SetResult struct {
	Key      string
	Previous string
	Updated  bool // false means it was a new key
}

// SetKey adds or updates a single key in the encrypted vault file.
// If overwrite is false and the key already exists, an error is returned.
func (v *Vault) SetKey(key, value string, overwrite bool) (*SetResult, error) {
	if key == "" {
		return nil, fmt.Errorf("key must not be empty")
	}

	// Decrypt current contents.
	entries, err := v.decryptEntries()
	if err != nil {
		return nil, fmt.Errorf("set: decrypt: %w", err)
	}

	result := &SetResult{Key: key}

	// Check for existing key.
	for _, e := range entries {
		if e.Key == key {
			if !overwrite {
				return nil, fmt.Errorf("key %q already exists; use --overwrite to replace it", key)
			}
			result.Previous = e.Value
			result.Updated = true
			break
		}
	}

	// Apply the change.
	found := false
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			found = true
			break
		}
	}
	if !found {
		entries = append(entries, env.Entry{Key: key, Value: value})
	}

	// Re-encrypt.
	if err := v.encryptEntries(entries); err != nil {
		return nil, fmt.Errorf("set: encrypt: %w", err)
	}

	return result, nil
}

// DeleteKey removes a key from the encrypted vault file.
// Returns an error if the key does not exist.
func (v *Vault) DeleteKey(key string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}

	entries, err := v.decryptEntries()
	if err != nil {
		return fmt.Errorf("delete: decrypt: %w", err)
	}

	newEntries := entries[:0]
	found := false
	for _, e := range entries {
		if e.Key == key {
			found = true
			continue
		}
		newEntries = append(newEntries, e)
	}
	if !found {
		return fmt.Errorf("key %q not found", key)
	}

	return v.encryptEntries(newEntries)
}
