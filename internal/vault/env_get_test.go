package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func setupGetVault(t *testing.T, content string) *Vault {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	cfg := DefaultConfig(dir)
	cfg.EnvFile = envFile
	return &Vault{cfg: cfg}
}

func TestGetKeyBasic(t *testing.T) {
	v := setupGetVault(t, "FOO=bar\nBAZ=qux\n")
	results, err := v.GetKey([]string{"FOO"}, GetOptions{ExactMatch: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Found || results[0].Value != "bar" {
		t.Errorf("expected FOO=bar, got %+v", results)
	}
}

func TestGetKeyNotFound(t *testing.T) {
	v := setupGetVault(t, "FOO=bar\n")
	results, err := v.GetKey([]string{"MISSING"}, GetOptions{ExactMatch: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Found {
		t.Errorf("expected key to be not found")
	}
}

func TestGetKeyCaseInsensitive(t *testing.T) {
	v := setupGetVault(t, "MY_SECRET=hello\n")
	results, err := v.GetKey([]string{"my_secret"}, GetOptions{ExactMatch: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Found || results[0].Value != "hello" {
		t.Errorf("case-insensitive lookup failed: %+v", results)
	}
	if results[0].Key != "MY_SECRET" {
		t.Errorf("expected stored key name MY_SECRET, got %s", results[0].Key)
	}
}

func TestGetKeyMaskValue(t *testing.T) {
	v := setupGetVault(t, "TOKEN=supersecret\n")
	results, err := v.GetKey([]string{"TOKEN"}, GetOptions{ExactMatch: true, MaskValue: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value == "supersecret" {
		t.Errorf("expected masked value, got plaintext")
	}
}

func TestGetKeyMultiple(t *testing.T) {
	v := setupGetVault(t, "A=1\nB=2\nC=3\n")
	results, err := v.GetKey([]string{"A", "C", "Z"}, GetOptions{ExactMatch: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Found || results[0].Value != "1" {
		t.Errorf("A lookup failed: %+v", results[0])
	}
	if !results[1].Found || results[1].Value != "3" {
		t.Errorf("C lookup failed: %+v", results[1])
	}
	if results[2].Found {
		t.Errorf("Z should not be found")
	}
}

func TestGetKeyEmptyKeys(t *testing.T) {
	v := setupGetVault(t, "FOO=bar\n")
	_, err := v.GetKey([]string{}, GetOptions{})
	if err == nil {
		t.Error("expected error for empty keys slice")
	}
}
