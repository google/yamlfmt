// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hotfix_test

import (
	"testing"

	"github.com/google/yamlfmt/internal/hotfix"
)

func TestStripCRBytes(t *testing.T) {
	crlfContent := []byte("two\r\nlines\r\n")
	lfContent, _ := hotfix.StripCRBytes(crlfContent)
	count := countByte(lfContent, '\r')
	if count != 0 {
		t.Fatalf("Found %d CR (decimal 13) bytes in %v", count, lfContent)
	}
}

func TestWriteCRLF(t *testing.T) {
	lfContent := []byte("two\nlines\n")
	crlfContent, _ := hotfix.WriteCRLFBytes(lfContent)
	countCR := countByte(crlfContent, '\r')
	countLF := countByte(crlfContent, '\n')
	if countCR != countLF {
		t.Fatalf("Found %d CR (decimal 13) and %d LF (decimal 10) bytes in %v", countCR, countLF, crlfContent)
	}
}

func countByte(content []byte, searchByte byte) int {
	count := 0
	for _, b := range content {
		if b == searchByte {
			count++
		}
	}
	return count
}
