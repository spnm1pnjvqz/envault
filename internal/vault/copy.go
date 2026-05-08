package vault

import (
	"fmt"

	"github.com/user/envault/internal/env"
)

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	Copied    []string
	Skipped   []string
	Overwrote []string
}

// CopyOptions configures the behaviour of CopyKeys.
type CopyOptions struct {
	// Overwrite existing keys in the destination vault.
	Overwrite bool
	// Keys is an optional allowlist; if empty all keys are copied.
	Keys []string
}

// CopyKeys copies secrets from srcPath to dstPath, both of which must be
// existing, unlocked .env files managed by envault.
func CopyKeys(srcPath, dstPath string, opts CopyOptions) (*CopyResult, error) {
	srcEntries, err := env.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("copy: read source: %w", err)
	}

	dstEntries, err := env.ReadFile(dstPath)
	if err != nil {
		return nil, fmt.Errorf("copy: read destination: %w", err)
	}

	allowlist := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowlist[k] = true
	}

	dstMap := env.ToMap(dstEntries)
	srcMap := env.ToMap(srcEntries)

	result := &CopyResult{}

	for k, v := range srcMap {
		if len(allowlist) > 0 && !allowlist[k] {
			continue
		}
		if _, exists := dstMap[k]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			result.Overwrote = append(result.Overwrote, k)
		} else {
			result.Copied = append(result.Copied, k)
		}
		dstMap[k] = v
	}

	merged := env.FromMap(dstMap)
	if err := env.WriteFile(dstPath, merged); err != nil {
		return nil, fmt.Errorf("copy: write destination: %w", err)
	}

	return result, nil
}
