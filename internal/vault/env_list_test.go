package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func setupListVault(t *testing.T, content string) *Vault {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	cfg := DefaultConfig()
	cfg.EnvFile = envFile
	return &Vault{config: cfg}
}

func TestListKeysBasic(t *testing.T) {
	v := setupListVault(t, "FOO=bar\nBAZ=qux\nALPHA=1\n")
	res, err := v.ListKeys(ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	// Default sort by key: ALPHA, BAZ, FOO
	if res[0].Key != "ALPHA" || res[1].Key != "BAZ" || res[2].Key != "FOO" {
		t.Errorf("unexpected order: %v", res)
	}
}

func TestListKeysFilterPrefix(t *testing.T) {
	v := setupListVault(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=envault\n")
	res, err := v.ListKeys(ListOptions{FilterPrefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	for _, r := range res {
		if r.Key != "DB_HOST" && r.Key != "DB_PORT" {
			t.Errorf("unexpected key: %s", r.Key)
		}
	}
}

func TestListKeysMaskValues(t *testing.T) {
	v := setupListVault(t, "SECRET=mysecret\n")
	res, err := v.ListKeys(ListOptions{MaskValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result")
	}
	if res[0].Value == "mysecret" {
		t.Error("expected value to be masked")
	}
}

func TestListKeysNoSort(t *testing.T) {
	v := setupListVault(t, "ZZZ=1\nAAA=2\nMMM=3\n")
	res, err := v.ListKeys(ListOptions{SortBy: "none"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Key != "ZZZ" {
		t.Errorf("expected first key ZZZ, got %s", res[0].Key)
	}
}

func TestListKeysEmpty(t *testing.T) {
	v := setupListVault(t, "")
	res, err := v.ListKeys(ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	}
}

func TestFormatList(t *testing.T) {
	results := []ListResult{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux", Tags: []string{"prod", "secret"}},
	}
	out := FormatList(results, false)
	if out == "" {
		t.Error("expected non-empty output")
	}
	if len(out) < 10 {
		t.Errorf("output too short: %q", out)
	}
}

func TestFormatListEmpty(t *testing.T) {
	out := FormatList(nil, false)
	if out != "(no keys found)" {
		t.Errorf("unexpected output: %q", out)
	}
}
