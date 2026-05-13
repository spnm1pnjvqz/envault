package vault

import (
	"fmt"
	"strings"

	"github.com/jwhittle933/envault/internal/env"
)

// GetOptions configures the behaviour of GetKey.
type GetOptions struct {
	// MaskValue replaces the secret value with asterisks before returning.
	MaskValue bool
	// ExactMatch requires the key to match exactly (case-sensitive).
	ExactMatch bool
}

// GetResult holds the result of a single key lookup.
type GetResult struct {
	Key   string
	Value string
	Found bool
}

// GetKey retrieves one or more keys from the decrypted .env file.
// If ExactMatch is false the lookup is case-insensitive.
func (v *Vault) GetKey(keys []string, opts GetOptions) ([]GetResult, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("at least one key must be specified")
	}

	entries, err := env.ReadFile(v.cfg.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("reading env file: %w", err)
	}

	keyMap := env.ToMap(entries)

	results := make([]GetResult, 0, len(keys))
	for _, k := range keys {
		result := GetResult{Key: k}

		if opts.ExactMatch {
			val, ok := keyMap[k]
			if ok {
				result.Found = true
				result.Value = val
			}
		} else {
			for mk, mv := range keyMap {
				if strings.EqualFold(mk, k) {
					result.Key = mk // normalise to stored key name
					result.Found = true
					result.Value = mv
					break
				}
			}
		}

		if result.Found && opts.MaskValue {
			result.Value = env.MaskValue(result.Value)
		}

		results = append(results, result)
	}

	return results, nil
}
