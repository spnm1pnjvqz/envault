package vault

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joeshaw/envault/internal/env"
)

// EditOptions configures the Edit operation.
type EditOptions struct {
	// Editor overrides the $EDITOR environment variable.
	Editor string
}

// EditResult describes what changed after an edit session.
type EditResult struct {
	Modified bool
	KeysAdded []string
	KeysRemoved []string
	KeysChanged []string
}

// EditEnvFile opens the decrypted .env file in the user's preferred editor,
// then re-encrypts it after the editor exits. If the file was not modified
// the vault is left untouched.
func (v *Vault) EditEnvFile(opts EditOptions) (*EditResult, error) {
	pairs, err := v.loadPlaintext()
	if err != nil {
		return nil, fmt.Errorf("edit: load plaintext: %w", err)
	}

	// Write a temp file for the editor.
	tmp, err := os.CreateTemp("", "envault-edit-*.env")
	if err != nil {
		return nil, fmt.Errorf("edit: create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := env.WriteFile(tmpPath, pairs); err != nil {
		return nil, fmt.Errorf("edit: write temp file: %w", err)
	}

	before, _ := os.ReadFile(tmpPath)

	editor := opts.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("edit: editor exited with error: %w", err)
	}

	after, _ := os.ReadFile(tmpPath)
	if string(before) == string(after) {
		return &EditResult{Modified: false}, nil
	}

	newPairs, err := env.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("edit: parse edited file: %w", err)
	}

	result := buildEditResult(pairs, newPairs)

	if err := v.Lock(env.WriteFile, newPairs); err != nil {
		return nil, fmt.Errorf("edit: re-encrypt: %w", err)
	}

	return result, nil
}

func buildEditResult(before, after []env.Pair) *EditResult {
	bm := env.ToMap(before)
	am := env.ToMap(after)
	r := &EditResult{Modified: true}
	for k, av := range am {
		if bv, ok := bm[k]; !ok {
			r.KeysAdded = append(r.KeysAdded, k)
		} else if bv != av {
			r.KeysChanged = append(r.KeysChanged, k)
		}
	}
	for k := range bm {
		if _, ok := am[k]; !ok {
			r.KeysRemoved = append(r.KeysRemoved, k)
		}
	}
	return r
}
