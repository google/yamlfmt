// Copyright 2024 Google LLC
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

package yamlfmt_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/collections"
	"github.com/google/yamlfmt/internal/tempfile"
)

func TestFilepathCollector(t *testing.T) {
	testCaseTable{
		{
			name: "finds direct paths",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml"},
				{FileName: "z.yml"},
			},
			includePatterns: testPatterns{
				{pattern: "x.yaml"},
				{pattern: "y.yaml"},
				{pattern: "z.yml"},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml": {},
				"y.yaml": {},
				"z.yml":  {},
			},
		},
		{
			name: "finds all in directory one layer",
			files: []tempfile.Path{
				{FileName: "a", IsDir: true},
				{FileName: "a/x.yaml"},
				{FileName: "a/y.yaml"},
				{FileName: "a/z.yml"},
			},
			includePatterns: testPatterns{
				{pattern: "a"},
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: collections.Set[string]{
				"a/x.yaml": {},
				"a/y.yaml": {},
				"a/z.yml":  {},
			},
		},
		{
			name: "finds direct path to subdirectory",
			files: []tempfile.Path{
				{FileName: "a", IsDir: true},
				{FileName: "a/x.yaml"},
				{FileName: "a/y.yaml"},
				{FileName: "a/z.yml"},
			},
			includePatterns: testPatterns{
				{pattern: "a/x.yaml"},
				{pattern: "a/z.yml"},
			},
			expectedFiles: collections.Set[string]{
				"a/x.yaml": {},
				"a/z.yml":  {},
			},
		},
		{
			name: "finds all in layered directories",
			files: []tempfile.Path{
				{FileName: "a", IsDir: true},
				{FileName: "a/b", IsDir: true},
				{FileName: "x.yml"},
				{FileName: "y.yml"},
				{FileName: "z.yaml"},
				{FileName: "a/x.yaml"},
				{FileName: "a/b/x.yaml"},
				{FileName: "a/b/y.yml"},
			},
			includePatterns: testPatterns{
				{pattern: ""}, // with the test this functionally means the whole temp dir
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: collections.Set[string]{
				"x.yml":      {},
				"y.yml":      {},
				"z.yaml":     {},
				"a/x.yaml":   {},
				"a/b/x.yaml": {},
				"a/b/y.yml":  {},
			},
		},
		{
			name: "exclude files",
			files: []tempfile.Path{
				{FileName: "a", IsDir: true},
				{FileName: "a/b", IsDir: true},
				{FileName: "x.yml"},
				{FileName: "y.yml"},
				{FileName: "z.yaml"},
				{FileName: "a/x.yaml"},
				{FileName: "a/b/x.yaml"},
				{FileName: "a/b/y.yml"},
			},
			includePatterns: testPatterns{
				{pattern: ""}, // with the test this functionally means the whole temp dir
			},
			excludePatterns: testPatterns{
				{pattern: "x.yml"},
				{pattern: "a/x.yaml"},
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: collections.Set[string]{
				"y.yml":      {},
				"z.yaml":     {},
				"a/b/x.yaml": {},
				"a/b/y.yml":  {},
			},
		},
		{
			name:            "exclude directory",
			changeToTempDir: true,
			files: []tempfile.Path{
				{FileName: "x.yml"},
				{FileName: "y.yml"},
				{FileName: "z.yaml"},

				{FileName: "a", IsDir: true},
				{FileName: "a/x.yaml"},

				{FileName: "a/b", IsDir: true},
				{FileName: "a/b/x.yaml"},
				{FileName: "a/b/y.yml"},
			},
			includePatterns: testPatterns{
				{pattern: ""}, // with the test this functionally means the whole temp dir
			},
			excludePatterns: testPatterns{
				{pattern: "a/b"},
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: collections.Set[string]{
				"x.yml":    {},
				"y.yml":    {},
				"z.yaml":   {},
				"a/x.yaml": {},
			},
		},
		{
			name: "don't get files with wrong extension",
			files: []tempfile.Path{
				{FileName: "x.yml"},
				{FileName: "y.yaml"},
				{FileName: "z.json"},
			},
			includePatterns: testPatterns{
				{pattern: ""}, // with the test this functionally means the whole temp dir
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: collections.Set[string]{
				"x.yml":  {},
				"y.yaml": {},
			},
		},
		{
			name: "multi-part extension",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml.gotmpl"},
				{FileName: "z.json"},
			},
			includePatterns: testPatterns{
				{pattern: ""}, // with the test this functionally means the whole temp dir
			},
			extensions: []string{
				"yaml",
				"yml",
				"yaml.gotmpl",
			},
			expectedFiles: collections.Set[string]{
				"x.yaml":        {},
				"y.yaml.gotmpl": {},
			},
		},
	}.runAll(t, useFilepathCollector)
}

func TestDoublestarCollectorBasic(t *testing.T) {
	testCaseTable{
		{
			name: "no excludes",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml"},
				{FileName: "z.yaml"},
			},
			includePatterns: testPatterns{
				{pattern: "**/*.yaml"},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml": {},
				"y.yaml": {},
				"z.yaml": {},
			},
		},
	}.runAll(t, useDoublestarCollector)
}

func TestDoublestarCollectorExcludeDirectory(t *testing.T) {
	testFiles := []tempfile.Path{
		{FileName: "x.yaml"},

		{FileName: "y", IsDir: true},
		{FileName: "y/y.yaml"},

		{FileName: "z", IsDir: true},
		{FileName: "z/z.yaml"},
		{FileName: "z/z1.yaml"},
		{FileName: "z/z2.yaml"},
	}

	testCaseTable{
		{
			name:  "exclude_directory/start with doublestar",
			files: testFiles,
			includePatterns: testPatterns{
				{pattern: "**/*.yaml"},
			},
			excludePatterns: testPatterns{
				{pattern: "**/z/**/*.yaml", stayRelative: true},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml":   {},
				"y/y.yaml": {},
			},
		},
		{
			name:            "exclude_directory/relative include and exclude",
			changeToTempDir: true,
			files:           testFiles,
			includePatterns: testPatterns{
				{pattern: "**/*.yaml", stayRelative: true},
			},
			excludePatterns: testPatterns{
				{pattern: "z/**/*.yaml", stayRelative: true},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml":   {},
				"y/y.yaml": {},
			},
		},
		{
			name:  "exclude_directory/absolute include and exclude",
			files: testFiles,
			includePatterns: testPatterns{
				{pattern: "**/*.yaml"},
			},
			excludePatterns: testPatterns{
				{pattern: "z/**/*.yaml"},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml":   {},
				"y/y.yaml": {},
			},
		},
		{
			name:            "exclude_directory/absolute include relative exclude",
			skip:            true,
			changeToTempDir: true,
			files:           testFiles,
			includePatterns: testPatterns{
				{pattern: "**/*.yaml"},
			},
			excludePatterns: testPatterns{
				{pattern: "z/**/*.yaml", stayRelative: true},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml":   {},
				"y/y.yaml": {},
			},
		},
		{
			name:            "exclude_directory/relative include absolute exclude",
			skip:            true,
			changeToTempDir: true,
			files:           testFiles,
			includePatterns: testPatterns{
				{pattern: "**/*.yaml", stayRelative: true},
			},
			excludePatterns: testPatterns{
				{pattern: "z/**/*.yaml"},
			},
			expectedFiles: collections.Set[string]{
				"x.yaml":   {},
				"y/y.yaml": {},
			},
		},
	}.runAll(t, useDoublestarCollector)
}

type testPatterns []struct {
	pattern      string
	stayRelative bool
}

func (tps testPatterns) allPatterns(path string) []string {
	result := make([]string, len(tps))
	for i := 0; i < len(tps); i++ {
		if tps[i].stayRelative {
			result[i] = tps[i].pattern
		} else {
			result[i] = fmt.Sprintf("%s/%s", path, tps[i].pattern)
		}
	}
	return result
}

// In some test scenarios we want to ignore whether a pattern is marked stayRelative
// and always treat them as relative by formatting the base path on them.
func (tps testPatterns) allPatternsForceAbsolute(path string) []string {
	result := make([]string, len(tps))
	for i := 0; i < len(tps); i++ {
		result[i] = fmt.Sprintf("%s/%s", path, tps[i].pattern)
	}
	return result
}

type testCase struct {
	name            string
	skip            bool
	changeToTempDir bool
	files           []tempfile.Path
	includePatterns testPatterns
	extensions      []string
	excludePatterns testPatterns
	expectedFiles   collections.Set[string]
}

func (tc testCase) run(t *testing.T, makeCollector makeCollectorFunc) {
	testStartDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	t.Run(tc.name, func(t *testing.T) {
		if tc.skip {
			t.Skip()
		}
		tempPath := t.TempDir()

		if tc.changeToTempDir {
			os.Chdir(tempPath)
		}

		for _, file := range tc.files {
			file.BasePath = tempPath
			if err := file.Create(); err != nil {
				t.Fatalf("Failed to create file")
			}
		}

		collector := makeCollector(tc, tempPath)
		paths, err := collector.CollectPaths()
		if err != nil {
			t.Fatalf("Test case failed: %v", err)
		}

		filesToFormat := collections.Set[string]{}
		for _, path := range paths {
			formatPath := path
			if strings.HasPrefix(formatPath, "/") {
				formatPath, err = filepath.Rel(tempPath, path)
				if err != nil {
					t.Fatalf("Path %s could not match to path %s", tempPath, path)
				}
			}
			filesToFormat.Add(formatPath)
		}
		if !filesToFormat.Equals(tc.expectedFiles) {
			t.Fatalf("Expected to receive paths %v\nbut got %v", tc.expectedFiles, filesToFormat)
		}
	})

	// Restore the starting directory if we changed in the test.
	if tc.changeToTempDir {
		os.Chdir(testStartDir)
	}
}

type testCaseTable []testCase

func (tcs testCaseTable) runAll(t *testing.T, makeCollector makeCollectorFunc) {
	for _, tc := range tcs {
		tc.run(t, makeCollector)
	}
}

type makeCollectorFunc func(tc testCase, path string) yamlfmt.PathCollector

func useFilepathCollector(tc testCase, path string) yamlfmt.PathCollector {
	return &yamlfmt.FilepathCollector{
		Include:    tc.includePatterns.allPatterns(path),
		Exclude:    tc.excludePatterns.allPatterns(path),
		Extensions: tc.extensions,
	}
}

func useDoublestarCollector(tc testCase, path string) yamlfmt.PathCollector {
	var includePatterns []string
	if tc.changeToTempDir {
		includePatterns = tc.includePatterns.allPatterns(path)
	} else {
		// If we didn't change to temp dir, disallow relative paths so we don't pick up
		// something confusing from the main working directory.
		includePatterns = tc.includePatterns.allPatternsForceAbsolute(path)
	}
	return &yamlfmt.DoublestarCollector{
		Include: includePatterns,
		Exclude: tc.excludePatterns.allPatterns(path),
	}
}

func TestPatternFile(t *testing.T) {
	t.Parallel()

	makePatterns := func(patterns ...string) []byte {
		var b bytes.Buffer

		fmt.Fprintln(&b, "# Comment followed by empty line")
		fmt.Fprintln(&b)
		for _, p := range patterns {
			fmt.Fprintln(&b, p)
		}

		return b.Bytes()
	}

	cases := []struct {
		name      string
		patterns  []byte
		haveFiles []string
		wantFiles []string
	}{
		{
			name:      "yaml and yml files",
			patterns:  makePatterns("*.yaml", "*.yml"),
			haveFiles: []string{"x.yaml", "y.yml", "README.md"},
			wantFiles: []string{"x.yaml", "y.yml"},
		},
		{
			name:      "ignore pattern",
			patterns:  makePatterns("*.yaml", "*.yml", "!test_input.yaml"),
			haveFiles: []string{"x.yaml", "y.yml", "test_input.yaml"},
			wantFiles: []string{"x.yaml", "y.yml"},
		},
		{
			name:      "descent into directories",
			patterns:  makePatterns("*.yaml", "*.yml"),
			haveFiles: []string{"a/x.yaml", "b/y.yml"},
			wantFiles: []string{"a/x.yaml", "b/y.yml"},
		},
		{
			name:      "exclude directories",
			patterns:  makePatterns("*.yaml", "*.yml", "!a/"),
			haveFiles: []string{"a/x.yaml", "b/y.yml"},
			wantFiles: []string{"b/y.yml"},
		},
		{
			name:      "matches are rooted at the working directory",
			patterns:  makePatterns("*.yaml", "!/x.yaml"),
			haveFiles: []string{"x.yaml", "a/x.yaml"},
			wantFiles: []string{"a/x.yaml"},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := make(fstest.MapFS)
			for _, f := range tc.haveFiles {
				fs[f] = &fstest.MapFile{
					Data: []byte("test"),
				}
			}

			patternFile := yamlfmt.NewPatternFileCollectorFS(bytes.NewReader(tc.patterns), fs)

			gotFiles, err := patternFile.CollectPaths()
			if err != nil {
				t.Fatal(err)
			}

			// Ignore the order of files in tc.wantFiles and gotFiles.
			opts := []cmp.Option{
				cmpopts.SortSlices(func(a, b string) bool { return a < b }),
			}

			if diff := cmp.Diff(tc.wantFiles, gotFiles, opts...); diff != "" {
				t.Errorf("PatternFile.CollectPaths() differs (-want/+got):\n%s", diff)
			}
		})
	}
}
