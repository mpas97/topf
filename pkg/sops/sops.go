// Copyright 2026 PostFinance AG
// SPDX-License-Identifier: MIT

// Package sops provides utilities for reading SOPS-encrypted files
package sops

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// ReadFileWithSOPS reads a file and automatically decrypts it with SOPS if encrypted.
// Returns the file contents, or nil if the file doesn't exist (not an error).
func ReadFileWithSOPS(path string) ([]byte, error) {
	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// not an error, just no file available
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Check if file is encrypted with SOPS
	// #nosec G204 required as long as we don't inline sops encryption
	statusCmd := exec.Command("sops", "filestatus", path)
	statusOutput, statusErr := statusCmd.Output()

	var status struct {
		Encrypted bool `json:"encrypted"`
	}

	isEncrypted := statusErr == nil &&
		json.Unmarshal(statusOutput, &status) == nil &&
		status.Encrypted

	if !isEncrypted {
		//nolint:gosec // files read through a variable in our control
		return os.ReadFile(path)
	}

	// File is encrypted: try to decrypt
	// #nosec G204 required as long as we don't inline sops encryption
	decryptCmd := exec.Command("sops", "decrypt", path)

	output, err := decryptCmd.Output()
	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("sops decryption failed: %s", string(exitErr.Stderr))
		}

		return nil, fmt.Errorf("failed to run sops decrypt: %w", err)
	}

	return output, nil
}
