package vault

import (
	"fmt"

	"github.com/nicholasgasior/envault/internal/env"
)

// DeleteOptions configures the behaviour of DeleteKeys.
type DeleteOptions struct {
	// Keys is the list of key names to remove.
	Keys []string
	// MustExist causes DeleteKeys to return an error if any key is not present.
	MustExist bool
}

// DeleteResult summarises what happened during a delete operation.
type DeleteResult struct {
	Deleted []string
	NotFound []string
}

// DeleteKeys removes one or more keys from the decrypted env file located at
// cfg.EnvFile. The file is read, modified in-memory and written back.
func (v *Vault) DeleteKeys(opts DeleteOptions) (*DeleteResult, error) {
	if len(opts.Keys) == 0 {
		return nil, fmt.Errorf("delete: no keys specified")
	}

	entries, err := env.ReadFile(v.cfg.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("delete: read env file: %w", err)
	}

	toDelete := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		toDelete[k] = struct{}{}
	}

	result := &DeleteResult{}
	present := make(map[string]struct{})
	for _, e := range entries {
		present[e.Key] = struct{}{}
	}

	for k := range toDelete {
		if _, ok := present[k]; ok {
			result.Deleted = append(result.Deleted, k)
		} else {
			result.NotFound = append(result.NotFound, k)
		}
	}

	if opts.MustExist && len(result.NotFound) > 0 {
		return result, fmt.Errorf("delete: keys not found: %v", result.NotFound)
	}

	filtered := entries[:0]
	for _, e := range entries {
		if _, remove := toDelete[e.Key]; !remove {
			filtered = append(filtered, e)
		}
	}

	if err := env.WriteFile(v.cfg.EnvFile, filtered); err != nil {
		return nil, fmt.Errorf("delete: write env file: %w", err)
	}

	return result, nil
}
