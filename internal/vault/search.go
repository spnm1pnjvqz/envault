package vault

import (
	"strings"

	"github.com/yourusername/envault/internal/env"
)

// SearchResult holds a matched key-value pair and optional context.
type SearchResult struct {
	Key   string
	Value string
	Masked bool
}

// SearchOptions controls how a search is performed.
type SearchOptions struct {
	// Query is the substring or prefix to match against keys.
	Query string
	// CaseSensitive controls whether matching is case-sensitive.
	CaseSensitive bool
	// MaskValues replaces values with asterisks in results.
	MaskValues bool
	// KeyOnly restricts matching to keys only (ignores values).
	KeyOnly bool
}

// Search decrypts the vault and returns entries whose keys (or values)
// match the given query string. It never writes to disk.
func (v *Vault) Search(opts SearchOptions) ([]SearchResult, error) {
	entries, err := v.decrypt()
	if err != nil {
		return nil, err
	}

	query := opts.Query
	if !opts.CaseSensitive {
		query = strings.ToLower(query)
	}

	var results []SearchResult
	for _, entry := range entries {
		key := entry.Key
		val := entry.Value

		cmpKey := key
		cmpVal := val
		if !opts.CaseSensitive {
			cmpKey = strings.ToLower(key)
			cmpVal = strings.ToLower(val)
		}

		keyMatch := strings.Contains(cmpKey, query)
		valMatch := !opts.KeyOnly && strings.Contains(cmpVal, query)

		if keyMatch || valMatch {
			displayVal := val
			if opts.MaskValues {
				displayVal = env.MaskValue(val)
			}
			results = append(results, SearchResult{
				Key:    key,
				Value:  displayVal,
				Masked: opts.MaskValues,
			})
		}
	}
	return results, nil
}
