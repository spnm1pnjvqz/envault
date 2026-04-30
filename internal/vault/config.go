package vault

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	DefaultConfigDir  = ".envault"
	DefaultConfigFile = "config.json"
	DefaultKeyFile    = "identity.age"
	DefaultEnvFile    = ".env"
	DefaultVaultFile  = ".env.age"
)

// Config holds the envault project configuration.
type Config struct {
	IdentityPath string `json:"identity_path"`
	EnvFile      string `json:"env_file"`
	VaultFile    string `json:"vault_file"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig(homeDir string) *Config {
	return &Config{
		IdentityPath: filepath.Join(homeDir, DefaultConfigDir, DefaultKeyFile),
		EnvFile:      DefaultEnvFile,
		VaultFile:    DefaultVaultFile,
	}
}

// SaveConfig writes the config as JSON to the given path.
func SaveConfig(path string, cfg *Config) error {
	if path == "" {
		return errors.New("config path must not be empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

// LoadConfig reads and parses a Config from the given path.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
