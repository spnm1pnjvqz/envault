package vault

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
)

const (
	// DefaultKeyFile is the default location for the age identity file.
	DefaultKeyFile = "~/.config/envault/key.txt"
)

// KeyPair holds an age recipient (public key) and identity (private key).
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// GenerateKeyPair creates a new age X25519 key pair.
func GenerateKeyPair() (*KeyPair, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		PublicKey:  identity.Recipient().String(),
		PrivateKey: identity.String(),
	}, nil
}

// SaveKeyPair writes the private key to a file, creating directories as needed.
func SaveKeyPair(kp *KeyPair, path string) error {
	if path == "" {
		return errors.New("key path must not be empty")
	}
	expanded := expandHome(path)
	if err := os.MkdirAll(filepath.Dir(expanded), 0700); err != nil {
		return err
	}
	content := "# envault age identity file\n# public key: " + kp.PublicKey + "\n" + kp.PrivateKey + "\n"
	return os.WriteFile(expanded, []byte(content), 0600)
}

// LoadPrivateKey reads an age identity (private key) from a file.
func LoadPrivateKey(path string) (string, error) {
	expanded := expandHome(path)
	data, err := os.ReadFile(expanded)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "AGE-SECRET-KEY-") {
			return line, nil
		}
	}
	return "", errors.New("no age secret key found in file")
}

// KeyFileExists reports whether a key file exists at the given path.
// It expands a leading ~ to the user's home directory before checking.
func KeyFileExists(path string) bool {
	expanded := expandHome(path)
	_, err := os.Stat(expanded)
	return err == nil
}

// expandHome replaces a leading ~ with the user's home directory.
func expandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}
