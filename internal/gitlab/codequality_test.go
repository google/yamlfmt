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
	"github.com/google/yamlfmt/internal/assert"
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
			assert.Equal(t, tc.wantOK, gotOK)
			if !gotOK {
				return
			}

			if tc.wantFingerprint != "" {
				assert.Equal(t, tc.wantFingerprint, got.Fingerprint)
			}

			data, err := json.Marshal(got)
			assert.NilErr(t, err)

			var gotUnmarshal gitlab.CodeQuality
			err = json.Unmarshal(data, &gotUnmarshal)
			assert.NilErr(t, err)

			if d := cmp.Diff(got, gotUnmarshal); d != "" {
				assert.EqualMsg(t, "", d, "json.Marshal() and json.Unmarshal() mismatch (-got +want):\n%s")
			}
		})
	}
}

func TestCodeQuality_DetectChangedLine(t *testing.T) {
	t.Parallel()

	testdataDir := "./testdata/changed_line"
	originalPath := filepath.Join(testdataDir, "original.yaml")
	formattedPath := filepath.Join(testdataDir, "formatted.yaml")

	original, err := os.ReadFile(originalPath)
	assert.NilErr(t, err)

	formatted, err := os.ReadFile(formattedPath)
	assert.NilErr(t, err)

	diff := yamlfmt.FileDiff{
		Path: "testdata/original.yaml",
		Diff: &yamlfmt.FormatDiff{
			Original:  original,
			Formatted: formatted,
		},
	}

	cq, ok := gitlab.NewCodeQuality(diff)
	assert.Assert(t, ok, "NewCodeQuality() returned false, expected true")

	assert.Assert(t, cq.Location.Lines != nil, "Location.Lines is nil")

	wantBeginLine := 6
	gotBeginLine := cq.Location.Lines.Begin
	assert.Equal(t, wantBeginLine, gotBeginLine)

	assert.Assert(t, cq.Location.Lines.End != nil, "Location.Lines.End is nil")

	wantEndLine := 8
	gotEndLine := *cq.Location.Lines.End
	assert.Equal(t, wantEndLine, gotEndLine)

	assert.Equal(t, diff.Path, cq.Location.Path)

	assert.Assert(t, cq.Description != "", "Description is empty")
	assert.Assert(t, cq.Name != "", "Name is empty")
	assert.Assert(t, cq.Fingerprint != "", "Fingerprint is empty")
	assert.Assert(t, cq.Severity != "", "Severity is empty")
}

func TestCodeQuality_DetectChangedLines_FromTestdata(t *testing.T) {
	t.Parallel()

	type tc struct {
		name      string
		dir       string
		wantOK    bool
		wantBegin int
		wantEnd   int
	}

	cases := []tc{
		{
			name:   "no lines changed",
			dir:    "no_lines_changed",
			wantOK: false,
		},
		{
			name:      "all lines changed",
			dir:       "all_lines_changed",
			wantOK:    true,
			wantBegin: 1,
			wantEnd:   2,
		},
		{
			name:      "single line changed",
			dir:       "single_line_changed",
			wantOK:    true,
			wantBegin: 2,
			wantEnd:   2,
		},
		{
			name:      "only the last line changed",
			dir:       "last_line_changed",
			wantOK:    true,
			wantBegin: 3,
			wantEnd:   3,
		},
		{
			name:      "only change is appending a last line",
			dir:       "append_last_line",
			wantOK:    true,
			wantBegin: 3,
			wantEnd:   3,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			testdataDir := filepath.Join("./testdata/changed_line", c.dir)
			originalPath := filepath.Join(testdataDir, "original.yaml")
			formattedPath := filepath.Join(testdataDir, "formatted.yaml")

			original, err := os.ReadFile(originalPath)
			assert.NilErr(t, err)

			formatted, err := os.ReadFile(formattedPath)
			assert.NilErr(t, err)

			diff := yamlfmt.FileDiff{
				Path: filepath.ToSlash(filepath.Join("testdata/changed_line", c.dir, "original.yaml")),
				Diff: &yamlfmt.FormatDiff{
					Original:  original,
					Formatted: formatted,
				},
			}

			cq, ok := gitlab.NewCodeQuality(diff)
			assert.Equal(t, c.wantOK, ok)
			if !ok {
				return
			}

			assert.Assert(t, cq.Location.Lines != nil, "Location.Lines is nil")
			assert.Equal(t, c.wantBegin, cq.Location.Lines.Begin)

			assert.Assert(t, cq.Location.Lines.End != nil, "Location.Lines.End is nil")
			assert.Equal(t, c.wantEnd, *cq.Location.Lines.End)

			assert.Equal(t, diff.Path, cq.Location.Path)
			assert.Assert(t, cq.Description != "", "Description is empty")
			assert.Assert(t, cq.Name != "", "Name is empty")
			assert.Assert(t, cq.Fingerprint != "", "Fingerprint is empty")
			assert.Assert(t, cq.Severity != "", "Severity is empty")
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
					Original:  []byte(tc.original),
					Formatted: []byte(tc.formatted),
				},
			}

			cq, ok := gitlab.NewCodeQuality(diff)
			assert.Assert(t, ok, "NewCodeQuality() returned false, expected true")

			assert.Assert(t, cq.Location.Lines != nil, "Location.Lines is nil")
			assert.Equal(t, tc.wantBegin, cq.Location.Lines.Begin)

			assert.Assert(t, cq.Location.Lines.End != nil, "Location.Lines.End is nil")
			assert.Equal(t, tc.wantEnd, *cq.Location.Lines.End)
		})
	}
}
