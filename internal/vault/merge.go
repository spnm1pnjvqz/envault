package vault

import (
	"fmt"

	"github.com/nicholasgasior/envault/internal/env"
)

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy string

const (
	// MergeStrategyOurs keeps the destination value on conflict.
	MergeStrategyOurs MergeStrategy = "ours"
	// MergeStrategyTheirs overwrites with the source value on conflict.
	MergeStrategyTheirs MergeStrategy = "theirs"
	// MergeStrategyError returns an error on any conflict.
	MergeStrategyError MergeStrategy = "error"
)

// MergeResult summarises the outcome of a merge operation.
type MergeResult struct {
	Added    []string
	Skipped  []string
	Overwritten []string
}

// MergeEnvFiles merges key/value pairs from src into dst.
// Existing keys in dst are handled according to the given strategy.
func MergeEnvFiles(dstPath, srcPath string, strategy MergeStrategy) (*MergeResult, error) {
	dstEntries, err := env.ReadFile(dstPath)
	if err != nil {
		return nil, fmt.Errorf("merge: read dst %q: %w", dstPath, err)
	}

	srcEntries, err := env.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("merge: read src %q: %w", srcPath, err)
	}

	dstMap := env.ToMap(dstEntries)
	result := &MergeResult{}

	for _, entry := range srcEntries {
		if entry.Comment || entry.Key == "" {
			continue
		}
		_, exists := dstMap[entry.Key]
		switch {
		case !exists:
			dstEntries = append(dstEntries, entry)
			dstMap[entry.Key] = entry.Value
			result.Added = append(result.Added, entry.Key)
		case strategy == MergeStrategyTheirs:
			for i, e := range dstEntries {
				if e.Key == entry.Key {
					dstEntries[i].Value = entry.Value
					break
				}
			}
			result.Overwritten = append(result.Overwritten, entry.Key)
		case strategy == MergeStrategyError:
			return nil, fmt.Errorf("merge: conflict on key %q", entry.Key)
		default: // MergeStrategyOurs
			result.Skipped = append(result.Skipped, entry.Key)
		}
	}

	if err := env.WriteFile(dstPath, dstEntries); err != nil {
		return nil, fmt.Errorf("merge: write dst %q: %w", dstPath, err)
	}

	return result, nil
}
