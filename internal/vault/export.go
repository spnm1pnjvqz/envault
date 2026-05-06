package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/envault/internal/env"
)

// ExportFormat represents the supported export formats.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
)

// ExportOptions configures the export operation.
type ExportOptions struct {
	Format    ExportFormat
	OutputPath string
	MaskSecrets bool
}

// Export decrypts the vault and writes its contents to a file or stdout
// in the requested format.
func (v *Vault) Export(privateKeyPath string, opts ExportOptions) error {
	entries, err := v.View(privateKeyPath)
	if err != nil {
		return fmt.Errorf("export: failed to decrypt vault: %w", err)
	}

	if opts.MaskSecrets {
		for i, e := range entries {
			entries[i].Value = env.MaskValue(e.Value)
		}
	}

	var data []byte
	switch opts.Format {
	case FormatJSON:
		data, err = exportJSON(entries)
	case FormatDotenv, "":
		data, err = exportDotenv(entries)
	default:
		return fmt.Errorf("export: unsupported format %q", opts.Format)
	}
	if err != nil {
		return fmt.Errorf("export: serialization failed: %w", err)
	}

	if opts.OutputPath == "" {
		_, err = os.Stdout.Write(data)
		return err
	}

	if err := os.WriteFile(opts.OutputPath, data, 0600); err != nil {
		return fmt.Errorf("export: write file: %w", err)
	}
	return nil
}

func exportDotenv(entries []env.Entry) ([]byte, error) {
	out := env.Serialize(entries)
	return []byte(out), nil
}

func exportJSON(entries []env.Entry) ([]byte, error) {
	type jsonEntry struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		Comment   string `json:"comment,omitempty"`
		ExportedAt string `json:"exported_at"`
	}

	now := time.Now().UTC().Format(time.RFC3339)
	records := make([]jsonEntry, 0, len(entries))
	for _, e := range entries {
		records = append(records, jsonEntry{
			Key:        e.Key,
			Value:      e.Value,
			Comment:    e.Comment,
			ExportedAt: now,
		})
	}
	return json.MarshalIndent(records, "", "  ")
}
