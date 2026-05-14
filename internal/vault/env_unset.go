package vault

import (
	"fmt"

	"github.com/subtlepseudonym/envault/internal/env"
)

// UnsetOptions controls behaviour of the UnsetKeys operation.
type UnsetOptions struct {
	// MustExist causes UnsetKeys to return an error if any key is not found.
	MustExist bool
}

// UnsetResult summarises the outcome of an UnsetKeys call.
type UnsetResult struct {
	Removed []string
	NotFound []string
}

// UnsetKeys removes one or more keys from the decrypted env file managed by v.
// The file is read, modified in memory, and written back atomically.
func (v *Vault) UnsetKeys(keys []string, opts UnsetOptions) (*UnsetResult, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	entries, err := env.ReadFile(v.cfg.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}

	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if k == "" {
			return nil, fmt.Errorf("key must not be empty")
		}
		keySet[k] = struct{}{}
	}

	result := &UnsetResult{}
	found := make(map[string]bool, len(keys))

	var kept []env.Entry
	for _, e := range entries {
		if _, remove := keySet[e.Key]; remove {
			found[e.Key] = true
			result.Removed = append(result.Removed, e.Key)
		} else {
			kept = append(kept, e)
		}
	}

	for _, k := range keys {
		if !found[k] {
			if opts.MustExist {
				return nil, fmt.Errorf("key not found: %s", k)
			}
			result.NotFound = append(result.NotFound, k)
		}
	}

	if err := env.WriteFile(v.cfg.EnvFile, kept); err != nil {
		return nil, fmt.Errorf("write env file: %w", err)
	}

	return result, nil
}
