// Copyright 2024 GitLab, Inc.
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

package gitlab_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/gitlab"
)

func TestCodeQuality(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		diff            yamlfmt.FileDiff
		wantOK          bool
		wantFingerprint string
	}{
		{
			name: "no diff",
			diff: yamlfmt.FileDiff{
				Path: "testcase/no_diff.yaml",
				Diff: &yamlfmt.FormatDiff{
					Original:  []byte("a: b"),
					Formatted: []byte("a: b"),
				},
			},
			wantOK: false,
		},
		{
			name: "with diff",
			diff: yamlfmt.FileDiff{
				Path: "testcase/with_diff.yaml",
				Diff: &yamlfmt.FormatDiff{
					Original:  []byte("a:   b"),
					Formatted: []byte("a: b"),
				},
			},
			wantOK: true,
			// SHA256 of diff.Diff.Original
			wantFingerprint: "05088f1c296b4fd999a1efe48e4addd0f962a8569afbacc84c44630d47f09330",
		},
	}

	for _, tc := range cases {
		// copy tc to avoid capturing an aliased loop variable in a Goroutine.
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, gotOK := gitlab.NewCodeQuality(tc.diff)
			if gotOK != tc.wantOK {
				t.Fatalf("NewCodeQuality() = (%#v, %v), want (*, %v)", got, gotOK, tc.wantOK)
			}
			if !gotOK {
				return
			}

			if tc.wantFingerprint != "" && tc.wantFingerprint != got.Fingerprint {
				t.Fatalf("NewCodeQuality().Fingerprint = %q, want %q", got.Fingerprint, tc.wantFingerprint)
			}

			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}

			var gotUnmarshal gitlab.CodeQuality
			if err := json.Unmarshal(data, &gotUnmarshal); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, gotUnmarshal); diff != "" {
				t.Errorf("json.Marshal() and json.Unmarshal() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestCodeQuality_DetectChangedLines_MultipleCases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		original  string
		formatted string
		wantBegin int
		wantEnd   int
	}{
		{
			name:      "single line change",
			original:  "a:   b",
			formatted: "a: b",
			wantBegin: 1,
			wantEnd:   1,
		},
		{
			name: "multiple consecutive lines",
			original: `line1
line2:   value
line3:  value
line4`,
			formatted: `line1
line2: value
line3: value
line4`,
			wantBegin: 2,
			wantEnd:   3,
		},
		{
			name: "non-consecutive changes",
			original: `line1
line2:   value
line3
line4:  value
line5`,
			formatted: `line1
line2: value
line3
line4: value
line5`,
			wantBegin: 2,
			wantEnd:   4,
		},
		{
			name: "change at beginning",
			original: `key:   value
line2
line3`,
			formatted: `key: value
line2
line3`,
			wantBegin: 1,
			wantEnd:   1,
		},
		{
			name: "change at end",
			original: `line1
line2
key:   value`,
			formatted: `line1
line2
key: value`,
			wantBegin: 3,
			wantEnd:   3,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diff := yamlfmt.FileDiff{
				Path: "test.yaml",
				Diff: &yamlfmt.FormatDiff{
					Original:  tc.original,
					Formatted: tc.formatted,
				},
			}

			cq, ok := gitlab.NewCodeQuality(diff)
			if !ok {
				t.Fatal("NewCodeQuality() returned false, expected true")
			}

			if cq.Location.Lines == nil {
				t.Fatal("Location.Lines is nil")
			}

			if cq.Location.Lines.Begin == nil {
				t.Fatal("Location.Lines.Begin is nil")
			}

			if cq.Location.Lines.End == nil {
				t.Fatal("Location.Lines.End is nil")
			}

			gotBegin := *cq.Location.Lines.Begin
			if gotBegin != tc.wantBegin {
				t.Errorf("Location.Lines.Begin = %d, want %d", gotBegin, tc.wantBegin)
			}

			gotEnd := *cq.Location.Lines.End
			if gotEnd != tc.wantEnd {
				t.Errorf("Location.Lines.End = %d, want %d", gotEnd, tc.wantEnd)
			}
		})
	}
}

func TestCodeQuality_DetectChangedLines(t *testing.T) {
	t.Parallel()

	testdataDir := "testdata/gitlab/changed_line"
	print(testdataDir)
	originalPath := filepath.Join(testdataDir, "original.yaml")
	formattedPath := filepath.Join(testdataDir, "formatted.yaml")

	original, err := os.ReadFile(originalPath)
	if err != nil {
		t.Fatalf("failed to read original file: %v", err)
	}

	formatted, err := os.ReadFile(formattedPath)
	if err != nil {
		t.Fatalf("failed to read formatted file: %v", err)
	}

	diff := yamlfmt.FileDiff{
		Path: "testdata/original.yaml",
		Diff: &yamlfmt.FormatDiff{
			Original:  string(original),
			Formatted: string(formatted),
		},
	}

	cq, ok := gitlab.NewCodeQuality(diff)
	if !ok {
		t.Fatal("NewCodeQuality() returned false, expected true")
	}

	if cq.Location.Lines == nil {
		t.Fatal("Location.Lines is nil")
	}

	if cq.Location.Lines.Begin == nil {
		t.Fatal("Location.Lines.Begin is nil")
	}

	if cq.Location.Lines.End == nil {
		t.Fatal("Location.Lines.End is nil")
	}

	wantBeginLine := 6
	gotBeginLine := *cq.Location.Lines.Begin

	if gotBeginLine != wantBeginLine {
		t.Errorf("Location.Lines.Begin = %d, want %d", gotBeginLine, wantBeginLine)
	}

	wantEndLine := 7
	gotEndLine := *cq.Location.Lines.End

	if gotEndLine != wantEndLine {
		t.Errorf("Location.Lines.End = %d, want %d", gotEndLine, wantEndLine)
	}

	if cq.Location.Path != diff.Path {
		t.Errorf("Location.Path = %q, want %q", cq.Location.Path, diff.Path)
	}

	if cq.Description == "" {
		t.Error("Description is empty")
	}

	if cq.Name == "" {
		t.Error("Name is empty")
	}

	if cq.Fingerprint == "" {
		t.Error("Fingerprint is empty")
	}

	if cq.Severity == "" {
		t.Error("Severity is empty")
	}
}
