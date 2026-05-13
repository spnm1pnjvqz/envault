package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func setupSetVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	cfg := DefaultConfig(dir)
	cfg.PublicKey = pub
	cfg.EnvFile = filepath.Join(dir, ".env")
	if err := os.WriteFile(cfg.EnvFile, []byte("EXISTING=hello\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	v, err := New(cfg, priv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := v.Lock(); err != nil {
		t.Fatalf("Lock: %v", err)
	}
	return v, dir
}

func TestSetKeyNew(t *testing.T) {
	v, _ := setupSetVault(t)
	res, err := v.SetKey("NEW_KEY", "world", false)
	if err != nil {
		t.Fatalf("SetKey: %v", err)
	}
	if res.Updated {
		t.Error("expected Updated=false for new key")
	}
	if res.Previous != "" {
		t.Errorf("expected empty Previous, got %q", res.Previous)
	}
}

func TestSetKeyOverwrite(t *testing.T) {
	v, _ := setupSetVault(t)
	res, err := v.SetKey("EXISTING", "newval", true)
	if err != nil {
		t.Fatalf("SetKey overwrite: %v", err)
	}
	if !res.Updated {
		t.Error("expected Updated=true when overwriting")
	}
	if res.Previous != "hello" {
		t.Errorf("expected Previous=%q, got %q", "hello", res.Previous)
	}
}

func TestSetKeyNoOverwriteConflict(t *testing.T) {
	v, _ := setupSetVault(t)
	_, err := v.SetKey("EXISTING", "other", false)
	if err == nil {
		t.Fatal("expected error when overwrite=false and key exists")
	}
}

func TestSetKeyEmptyKey(t *testing.T) {
	v, _ := setupSetVault(t)
	_, err := v.SetKey("", "val", false)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestDeleteKey(t *testing.T) {
	v, _ := setupSetVault(t)
	if err := v.DeleteKey("EXISTING"); err != nil {
		t.Fatalf("DeleteKey: %v", err)
	}
}

func TestDeleteKeyNotFound(t *testing.T) {
	v, _ := setupSetVault(t)
	err := v.DeleteKey("DOES_NOT_EXIST")
	if err == nil {
		t.Fatal("expected error deleting non-existent key")
	}
}

func TestDeleteKeyEmptyKey(t *testing.T) {
	v, _ := setupSetVault(t)
	if err := v.DeleteKey(""); err == nil {
		t.Fatal("expected error for empty key")
	}
}
