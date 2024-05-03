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
	"path/filepath"
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/collections"
	"github.com/google/yamlfmt/internal/tempfile"
)

const testdataBase = "testdata/content_analyzer"

func TestBasicContentAnalyzer(t *testing.T) {
	testCases := []struct {
		name             string
		testdataDir      string
		excludePatterns  []string
		expectedPaths    collections.Set[string]
		expectedExcluded collections.Set[string]
	}{
		{
			name:            "has ignore metadata",
			testdataDir:     "has_ignore",
			excludePatterns: []string{},
			expectedPaths: collections.Set[string]{
				"y.yaml": {},
			},
			expectedExcluded: collections.Set[string]{
				"x.yaml": {},
			},
		},
		{
			name:        "matches regex pattern",
			testdataDir: "regex_ignore",
			excludePatterns: []string{
				".*generated by.*",
			},
			expectedPaths: collections.Set[string]{
				"y.yaml": {},
			},
			expectedExcluded: collections.Set[string]{
				"x.yaml": {},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempPath := t.TempDir()
			testdataDir := filepath.Join(testdataBase, tc.testdataDir)
			paths, err := tempfile.ReplicateDirectory(testdataDir, tempPath)
			if err != nil {
				t.Fatalf("could not replicate testdata directory %s: %v", tc.testdataDir, err)
			}
			err = paths.CreateAll()
			if err != nil {
				t.Fatalf("could not create full test directory: %v", err)
			}
			contentAnalyzer, err := yamlfmt.NewBasicContentAnalyzer(tc.excludePatterns)
			if err != nil {
				t.Fatalf("could not create content analyzer: %v", err)
			}
			collector := &yamlfmt.FilepathCollector{
				Include:    []string{tempPath},
				Exclude:    []string{},
				Extensions: []string{"yaml", "yml"},
			}
			collectedPaths, err := collector.CollectPaths()
			if err != nil {
				t.Fatalf("CollectPaths failed: %v", err)
			}
			resultPaths, excludedPaths, err := contentAnalyzer.ExcludePathsByContent(collectedPaths)
			if err != nil {
				t.Fatalf("expected content analyzer to work, got error: %v", err)
			}
			resultPathsTrimmed, err := pathsTempdirTrimmed(resultPaths, tempPath)
			if err != nil {
				t.Fatalf("expected trimming tempdir from result not to have error: %v", err)
			}
			if !tc.expectedPaths.Equals(collections.SliceToSet(resultPathsTrimmed)) {
				t.Fatalf("expected files:\n%v\ngot:\n%v", tc.expectedPaths, resultPaths)
			}
			excludePathsTrimmed, err := pathsTempdirTrimmed(excludedPaths, tempPath)
			if err != nil {
				t.Fatalf("expected trimming tempdir from excluded not to have error: %v", err)
			}
			if !tc.expectedExcluded.Equals(collections.SliceToSet(excludePathsTrimmed)) {
				t.Fatalf("expected exclusions:\n%v\ngot:\n%v", tc.expectedExcluded, excludedPaths)
			}
		})
	}
}

func TestBadNewContentAnalyzer(t *testing.T) {
	// Illegal because no closing )
	badPattern := "%^3412098(]fj"
	_, err := yamlfmt.NewBasicContentAnalyzer([]string{badPattern})
	if err == nil {
		t.Fatalf("expected there to be an error")
	}
}

func pathsTempdirTrimmed(paths []string, tempDir string) ([]string, error) {
	trimmedPaths := []string{}
	for _, path := range paths {
		trimmedPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return nil, err
		}
		trimmedPaths = append(trimmedPaths, trimmedPath)
	}
	return trimmedPaths, nil
}
