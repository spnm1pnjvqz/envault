package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTagVault(t *testing.T, content string) *Vault {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	cfg := DefaultConfig(dir)
	cfg.EnvFile = envFile
	v, err := New(cfg)
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	return v
}

func TestSetAndGetTags(t *testing.T) {
	v := setupTagVault(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	if err := SetTags(v, "DB_HOST", []string{"infra", "database"}); err != nil {
		t.Fatalf("SetTags: %v", err)
	}

	tags, err := GetTags(v, "DB_HOST")
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d: %v", len(tags), tags)
	}
}

func TestSetTagsSorted(t *testing.T) {
	v := setupTagVault(t, "API_KEY=secret\n")
	if err := SetTags(v, "API_KEY", []string{"zebra", "alpha"}); err != nil {
		t.Fatalf("SetTags: %v", err)
	}
	tags, _ := GetTags(v, "API_KEY")
	if tags[0] != "alpha" || tags[1] != "zebra" {
		t.Fatalf("expected sorted tags, got %v", tags)
	}
}

func TestClearTags(t *testing.T) {
	v := setupTagVault(t, "TOKEN=abc\n")
	_ = SetTags(v, "TOKEN", []string{"secret"})
	if err := SetTags(v, "TOKEN", []string{}); err != nil {
		t.Fatalf("ClearTags: %v", err)
	}
	tags, _ := GetTags(v, "TOKEN")
	if len(tags) != 0 {
		t.Fatalf("expected no tags, got %v", tags)
	}
}

func TestListByTag(t *testing.T) {
	v := setupTagVault(t, "DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=secret\n")
	_ = SetTags(v, "DB_HOST", []string{"database"})
	_ = SetTags(v, "DB_PORT", []string{"database"})
	_ = SetTags(v, "API_KEY", []string{"auth"})

	results, err := ListByTag(v, "database")
	if err != nil {
		t.Fatalf("ListByTag: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestListByTagNoMatch(t *testing.T) {
	v := setupTagVault(t, "FOO=bar\n")
	results, err := ListByTag(v, "nonexistent")
	if err != nil {
		t.Fatalf("ListByTag: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSetTagsKeyNotFound(t *testing.T) {
	v := setupTagVault(t, "FOO=bar\n")
	err := SetTags(v, "MISSING", []string{"tag"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
