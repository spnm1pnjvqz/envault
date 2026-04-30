package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}
	if !strings.HasPrefix(kp.PublicKey, "age1") {
		t.Errorf("expected public key to start with 'age1', got %q", kp.PublicKey)
	}
	if !strings.HasPrefix(kp.PrivateKey, "AGE-SECRET-KEY-") {
		t.Errorf("expected private key to start with 'AGE-SECRET-KEY-', got %q", kp.PrivateKey)
	}
}

func TestSaveAndLoadKeyPair(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "subdir", "key.txt")

	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	if err := SaveKeyPair(kp, keyPath); err != nil {
		t.Fatalf("SaveKeyPair() error: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}

	loaded, err := LoadPrivateKey(keyPath)
	if err != nil {
		t.Fatalf("LoadPrivateKey() error: %v", err)
	}
	if loaded != kp.PrivateKey {
		t.Errorf("loaded key %q does not match original %q", loaded, kp.PrivateKey)
	}
}

func TestSaveKeyPairEmptyPath(t *testing.T) {
	kp := &KeyPair{PublicKey: "age1test", PrivateKey: "AGE-SECRET-KEY-TEST"}
	if err := SaveKeyPair(kp, ""); err == nil {
		t.Error("expected error for empty path, got nil")
	}
}

func TestLoadPrivateKeyNotFound(t *testing.T) {
	_, err := LoadPrivateKey("/nonexistent/path/key.txt")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadPrivateKeyNoValidKey(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(keyPath, []byte("# just a comment\nno key here\n"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}
	_, err := LoadPrivateKey(keyPath)
	if err == nil {
		t.Error("expected error when no secret key found, got nil")
	}
}
