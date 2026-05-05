package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return buf.String(), err
}

func TestRootCommandHelp(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "envault",
		Short: "Lightweight secrets manager using age encryption",
	}
	cmd.AddCommand(initCmd, lockCmd, unlockCmd, viewCmd)

	out, err := executeCommand(cmd, "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected help output, got empty string")
	}
}

func TestLockCommandMissingConfig(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	cmd := &cobra.Command{Use: "envault"}
	cmd.AddCommand(lockCmd)
	_, err := executeCommand(cmd, "lock")
	if err == nil {
		t.Error("expected error when config is missing")
	}
}

func TestUnlockCommandMissingConfig(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	cmd := &cobra.Command{Use: "envault"}
	cmd.AddCommand(unlockCmd)
	_, err := executeCommand(cmd, "unlock")
	if err == nil {
		t.Error("expected error when config is missing")
	}
}

func TestViewCommandMissingConfig(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	cmd := &cobra.Command{Use: "envault"}
	cmd.AddCommand(viewCmd)
	_, err := executeCommand(cmd, "view")
	if err == nil {
		t.Error("expected error when config is missing")
	}
}

func TestConfigFileNameConvention(t *testing.T) {
	expected := ".envault.toml"
	if filepath.Base(expected) != ".envault.toml" {
		t.Errorf("unexpected config filename: %s", expected)
	}
}
