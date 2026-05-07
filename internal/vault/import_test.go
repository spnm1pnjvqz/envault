package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestImportDotenv(t *testing.T) {
	dir := t.TempDir()
	v, _ := setupImportVault(t, dir)

	src := filepath.Join(dir, "import.env")
	_ = os.WriteFile(src, []byte("NEW_KEY=hello\nDB_HOST=newhost\n"), 0600)

	result, err := v.Import(src, ImportFormatDotenv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Added) != 1 || result.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY added, got %v", result.Added)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST skipped, got %v", result.Skipped)
	}
}

func TestImportDotenvOverwrite(t *testing.T) {
	dir := t.TempDir()
	v, _ := setupImportVault(t, dir)

	src := filepath.Join(dir, "import.env")
	_ = os.WriteFile(src, []byte("DB_HOST=newhost\n"), 0600)

	result, err := v.Import(src, ImportFormatDotenv, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Updated) != 1 || result.Updated[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST updated, got %v", result.Updated)
	}
}

func TestImportJSON(t *testing.T) {
	dir := t.TempDir()
	v, _ := setupImportVault(t, dir)

	m := map[string]string{"API_KEY": "secret", "DB_HOST": "ignored"}
	data, _ := json.Marshal(m)
	src := filepath.Join(dir, "import.json")
	_ = os.WriteFile(src, data, 0600)

	result, err := v.Import(src, ImportFormatJSON, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Added) != 1 || result.Added[0] != "API_KEY" {
		t.Errorf("expected API_KEY added, got %v", result.Added)
	}
}

func TestImportUnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	v, _ := setupImportVault(t, dir)
	src := filepath.Join(dir, "import.env")
	_ = os.WriteFile(src, []byte(""), 0600)

	_, err := v.Import(src, "yaml", false)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestDetectFormat(t *testing.T) {
	cases := []struct {
		path   string
		want   ImportFormat
		wantErr bool
	}{
		{".env", ImportFormatDotenv, false},
		{"secrets.env", ImportFormatDotenv, false},
		{"secrets.json", ImportFormatJSON, false},
		{"secrets.yaml", "", true},
	}
	for _, tc := range cases {
		got, err := DetectFormat(tc.path)
		if tc.wantErr && err == nil {
			t.Errorf("DetectFormat(%q): expected error", tc.path)
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("DetectFormat(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestImportFileNotFound(t *testing.T) {
	dir := t.TempDir()
	v, _ := setupImportVault(t, dir)
	_, err := v.Import(filepath.Join(dir, "nonexistent.env"), ImportFormatDotenv, false)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func setupImportVault(t *testing.T, dir string) (*Vault, *Config) {
	t.Helper()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("DB_HOST=localhost\nDB_PORT=5432\n"), 0600)
	cfg := &Config{EnvFile: envFile}
	v := &Vault{config: cfg}
	return v, cfg
}
