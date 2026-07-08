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

// Package linescan provides a line scanner that is not subject to bufio's
// default token size limit.
package linescan

import (
	"bufio"
	"bytes"
)

// NewScanner returns a bufio.Scanner over content that can read a line of any length.
func NewScanner(content []byte) *bufio.Scanner {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(nil, len(content)+1)
	return scanner
}
