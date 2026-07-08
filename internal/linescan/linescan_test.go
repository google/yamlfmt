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

package linescan_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/google/yamlfmt/internal/linescan"
)

// A line longer than bufio's default token limit must scan without error,
// which a plain bufio.NewScanner cannot do.
func TestNewScannerLongLine(t *testing.T) {
	long := strings.Repeat("x", 10*bufio.MaxScanTokenSize) // 640KiB, one line
	content := []byte("a\n" + long + "\nb\n")

	var got []string
	scanner := linescan.NewScanner(content)
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	want := []string{"a", long, "b"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d mismatch: len(got)=%d len(want)=%d", i, len(got[i]), len(want[i]))
		}
	}
}

// Empty content is a valid (zero-line) scan, not an error.
func TestNewScannerEmpty(t *testing.T) {
	scanner := linescan.NewScanner(nil)
	for scanner.Scan() {
		t.Fatalf("empty content should yield no lines, got %q", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
}
