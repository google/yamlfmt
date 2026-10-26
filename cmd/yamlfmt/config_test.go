// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	testCases := []struct {
		name        string
		path        string
		expectedErr error
	}{
		{
			name:        "regular file",
			path:        filePath,
			expectedErr: nil,
		},
		{
			name:        "file does not exist",
			path:        filepath.Join(tempDir, ".yamlfmt"),
			expectedErr: errConfPathNotExist,
		},
		{
			name:        "path is a directory",
			path:        tempDir,
			expectedErr: errConfPathIsDir,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePath(tc.path)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestValidatePathUnderRegularFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Stat fails here with an error that is not "not exist" on most
	// platforms (ENOTDIR on Unix). It must be returned, not panic.
	err := validatePath(filepath.Join(filePath, ".yamlfmt"))
	if err == nil {
		t.Fatal("expected error for path under a regular file, got nil")
	}
	var pathErr *configPathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected a configPathError, got: %T %v", err, err)
	}
}
