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

package engine

import (
	"testing"

	"github.com/google/yamlfmt"
)

func TestEngineOutputSingleLineString(t *testing.T) {
	makeDiff := func(path, original, formatted string) *yamlfmt.FileDiff {
		return &yamlfmt.FileDiff{
			Path: path,
			Diff: &yamlfmt.FormatDiff{
				Original:  []byte(original),
				Formatted: []byte(formatted),
			},
		}
	}

	tests := []struct {
		name  string
		files yamlfmt.FileDiffs
		want  string
	}{
		{
			name: "only files with formatting differences are reported",
			files: yamlfmt.FileDiffs{
				"clean.yaml": makeDiff("clean.yaml", "a:\n  b: 1\n", "a:\n  b: 1\n"),
				"bad.yaml":   makeDiff("bad.yaml", "a:\n b: 1\n", "a:\n  b: 1\n"),
			},
			want: "bad.yaml: formatting difference found\n",
		},
		{
			name: "reported files are sorted by path",
			files: yamlfmt.FileDiffs{
				"z.yaml": makeDiff("z.yaml", "a:\n b: 1\n", "a:\n  b: 1\n"),
				"a.yaml": makeDiff("a.yaml", "a:\n b: 1\n", "a:\n  b: 1\n"),
				"m.yaml": makeDiff("m.yaml", "a:\n b: 1\n", "a:\n  b: 1\n"),
			},
			want: "a.yaml: formatting difference found\n" +
				"m.yaml: formatting difference found\n" +
				"z.yaml: formatting difference found\n",
		},
		{
			name: "no changed files reports nothing",
			files: yamlfmt.FileDiffs{
				"clean.yaml": makeDiff("clean.yaml", "a:\n  b: 1\n", "a:\n  b: 1\n"),
			},
			want: "",
		},
		{
			name:  "empty diff map reports nothing",
			files: yamlfmt.FileDiffs{},
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := (engineOutputSingleLine{Files: tc.files}).String()
			if got != tc.want {
				t.Fatalf("expected output %q, got %q", tc.want, got)
			}
		})
	}
}
