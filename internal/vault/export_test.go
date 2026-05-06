package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envault/internal/env"
)

func TestExportDotenv(t *testing.T) {
	pub, priv, dir := generateTestKeys(t)
	envFile := filepath.Join(dir, ".env")
	encFile := filepath.Join(dir, ".env.age")

	entries := []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
	if err := env.WriteFile(envFile, entries); err != nil {
		t.Fatal(err)
	}

	v := &Vault{Config: &Config{EnvFile: envFile, EncryptedFile: encFile, PublicKey: pub}}
	if err := v.Lock(); err != nil {
		t.Fatal(err)
	}

	outFile := filepath.Join(dir, "export.env")
	opts := ExportOptions{Format: FormatDotenv, OutputPath: outFile}
	if err := v.Export(priv, opts); err != nil {
		t.Fatalf("Export dotenv: %v", err)
	}

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got:\n%s", got)
	}
}

func TestExportJSON(t *testing.T) {
	pub, priv, dir := generateTestKeys(t)
	envFile := filepath.Join(dir, ".env")
	encFile := filepath.Join(dir, ".env.age")

	entries := []env.Entry{
		{Key: "SECRET", Value: "supersecret"},
	}
	if err := env.WriteFile(envFile, entries); err != nil {
		t.Fatal(err)
	}

	v := &Vault{Config: &Config{EnvFile: envFile, EncryptedFile: encFile, PublicKey: pub}}
	if err := v.Lock(); err != nil {
		t.Fatal(err)
	}

	outFile := filepath.Join(dir, "export.json")
	opts := ExportOptions{Format: FormatJSON, OutputPath: outFile}
	if err := v.Export(priv, opts); err != nil {
		t.Fatalf("Export json: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	var records []map[string]string
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 || records[0]["key"] != "SECRET" {
		t.Errorf("unexpected records: %v", records)
	}
}

func TestExportMaskSecrets(t *testing.T) {
	pub, priv, dir := generateTestKeys(t)
	envFile := filepath.Join(dir, ".env")
	encFile := filepath.Join(dir, ".env.age")

	entries := []env.Entry{{Key: "TOKEN", Value: "abc123xyz"}}
	if err := env.WriteFile(envFile, entries); err != nil {
		t.Fatal(err)
	}

	v := &Vault{Config: &Config{EnvFile: envFile, EncryptedFile: encFile, PublicKey: pub}}
	if err := v.Lock(); err != nil {
		t.Fatal(err)
	}

	outFile := filepath.Join(dir, "masked.env")
	opts := ExportOptions{Format: FormatDotenv, OutputPath: outFile, MaskSecrets: true}
	if err := v.Export(priv, opts); err != nil {
		t.Fatalf("Export masked: %v", err)
	}

	got, _ := os.ReadFile(outFile)
	if strings.Contains(string(got), "abc123xyz") {
		t.Errorf("expected value to be masked, got: %s", got)
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	pub, priv, dir := generateTestKeys(t)
	envFile := filepath.Join(dir, ".env")
	encFile := filepath.Join(dir, ".env.age")

	entries := []env.Entry{{Key: "X", Value: "1"}}
	_ = env.WriteFile(envFile, entries)
	v := &Vault{Config: &Config{EnvFile: envFile, EncryptedFile: encFile, PublicKey: pub}}
	_ = v.Lock()

	opts := ExportOptions{Format: "xml"}
	if err := v.Export(priv, opts); err == nil {
		t.Error("expected error for unsupported format")
	}
}
