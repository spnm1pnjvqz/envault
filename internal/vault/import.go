package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholasgasior/envault/internal/env"
)

// ImportFormat represents a supported import format.
type ImportFormat string

const (
	ImportFormatDotenv ImportFormat = "dotenv"
	ImportFormatJSON   ImportFormat = "json"
)

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Added    []string
	Updated  []string
	Skipped  []string
}

// Import reads secrets from an external file and merges them into the vault.
// If overwrite is false, existing keys are skipped.
func (v *Vault) Import(srcPath string, format ImportFormat, overwrite bool) (*ImportResult, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("reading import file: %w", err)
	}

	var incoming []env.Entry
	switch format {
	case ImportFormatDotenv:
		incoming, err = importDotenv(data)
	case ImportFormatJSON:
		incoming, err = importJSON(data)
	default:
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("parsing import file: %w", err)
	}

	existing, err := env.ReadFile(v.config.EnvFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading env file: %w", err)
	}

	existingMap := env.ToMap(existing)
	result := &ImportResult{}

	for _, entry := range incoming {
		if _, exists := existingMap[entry.Key]; exists {
			if !overwrite {
				result.Skipped = append(result.Skipped, entry.Key)
				continue
			}
			result.Updated = append(result.Updated, entry.Key)
		} else {
			result.Added = append(result.Added, entry.Key)
		}
		existingMap[entry.Key] = entry.Value
	}

	merged := env.FromMap(existingMap)
	if err := env.WriteFile(v.config.EnvFile, merged); err != nil {
		return nil, fmt.Errorf("writing env file: %w", err)
	}

	return result, nil
}

// DetectFormat attempts to detect the import format from the file extension.
func DetectFormat(path string) (ImportFormat, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".env", "":
		return ImportFormatDotenv, nil
	case ".json":
		return ImportFormatJSON, nil
	default:
		return "", fmt.Errorf("cannot detect format from extension %q", ext)
	}
}

func importDotenv(data []byte) ([]env.Entry, error) {
	return env.Parse(string(data))
}

func importJSON(data []byte) ([]env.Entry, error) {
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	entries := make([]env.Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, env.Entry{Key: k, Value: v})
	}
	return entries, nil
}
