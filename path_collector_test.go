package yamlfmt_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/tempfile"
)

func TestCollectPaths(t *testing.T) {
	testCases := []struct {
		name            string
		files           []tempfile.Path
		includePatterns []string
		excludePatterns []string
		extensions      []string
		expectedFiles   map[string]struct{}
	}{
		{
			name: "finds direct paths",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml"},
				{FileName: "z.yml"},
			},
			includePatterns: []string{
				"x.yaml",
				"y.yaml",
				"z.yml",
			},
			expectedFiles: map[string]struct{}{
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
			includePatterns: []string{
				"a",
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: map[string]struct{}{
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
			includePatterns: []string{
				"a/x.yaml",
				"a/z.yml",
			},
			expectedFiles: map[string]struct{}{
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
			includePatterns: []string{
				"", // with the test this functionally means the whole temp dir
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: map[string]struct{}{
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
			includePatterns: []string{
				"", // with the test this functionally means the whole temp dir
			},
			excludePatterns: []string{
				"x.yml",
				"a/x.yaml",
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: map[string]struct{}{
				"y.yml":      {},
				"z.yaml":     {},
				"a/b/x.yaml": {},
				"a/b/y.yml":  {},
			},
		},
		{
			name: "exclude directory",
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
			includePatterns: []string{
				"", // with the test this functionally means the whole temp dir
			},
			excludePatterns: []string{
				"a/b",
			},
			extensions: []string{
				"yaml",
				"yml",
			},
			expectedFiles: map[string]struct{}{
				"x.yml":    {},
				"y.yml":    {},
				"z.yaml":   {},
				"a/x.yaml": {},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempPath := t.TempDir()

			for _, file := range tc.files {
				file.BasePath = tempPath
				if err := file.Create(); err != nil {
					t.Fatalf("Failed to create file %s: %v", file.BasePath, err)
				}
			}

			collector := &yamlfmt.FilepathCollector{
				Include:    formatTempPaths(tempPath, tc.includePatterns),
				Exclude:    formatTempPaths(tempPath, tc.excludePatterns),
				Extensions: tc.extensions,
			}

			paths, err := collector.CollectPaths()
			if err != nil {
				t.Fatalf("CollectPaths failed: %v", err)
			}
			if len(paths) != len(tc.expectedFiles) {
				t.Fatalf("Got %d paths but expected %d", len(paths), len(tc.expectedFiles))
			}

			filesToFormat := map[string]struct{}{}
			for _, path := range paths {
				formatPath, err := filepath.Rel(tempPath, path)
				if err != nil {
					t.Fatalf("Path %s could match to path %s", tempPath, path)
				}
				filesToFormat[formatPath] = struct{}{}
			}
			if !reflect.DeepEqual(filesToFormat, tc.expectedFiles) {
				t.Fatalf("Expected to receive paths %v but got %v", tc.expectedFiles, filesToFormat)
			}
		})
	}
}

func TestDoublestarCollectPaths(t *testing.T) {
	t.Skip()
	testCases := []struct {
		name            string
		files           []tempfile.Path
		includePatterns []string
		excludePatterns []string
		expectedFiles   map[string]struct{}
	}{
		{
			name: "no excludes",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml"},
				{FileName: "z.yaml"},
			},
			expectedFiles: map[string]struct{}{
				"x.yaml": {},
				"y.yaml": {},
				"z.yaml": {},
			},
		},
		{
			name: "does not include directories",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y", IsDir: true},
			},
			expectedFiles: map[string]struct{}{
				"x.yaml": {},
			},
		},
		{
			name: "only include what is asked",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml"},
			},
			includePatterns: []string{
				"y.yaml",
			},
			expectedFiles: map[string]struct{}{
				"y.yaml": {},
			},
		},
		{
			name: "exclude what is asked",
			files: []tempfile.Path{
				{FileName: "x.yaml"},
				{FileName: "y.yaml"},
			},
			excludePatterns: []string{
				"y.yaml",
			},
			expectedFiles: map[string]struct{}{
				"x.yaml": {},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempPath := t.TempDir()

			for _, file := range tc.files {
				file.BasePath = tempPath
				if err := file.Create(); err != nil {
					t.Fatalf("Failed to create file")
				}
			}

			collector := &yamlfmt.DoublestarCollector{
				Include: formatTempPaths(tempPath, tc.includePatterns),
				Exclude: formatTempPaths(tempPath, tc.excludePatterns),
			}
			if len(collector.Include) == 0 {
				collector.Include = []string{fmt.Sprintf("%s/**", tempPath)}
			}
			if len(collector.Exclude) == 0 {
				collector.Exclude = []string{}
			}

			paths, err := collector.CollectPaths()
			if err != nil {
				t.Fatalf("CollectDoublestarPathsToFormat failed: %v", err)
			}

			filesToFormat := map[string]struct{}{}
			for _, path := range paths {
				_, filename := filepath.Split(path)
				filesToFormat[filename] = struct{}{}
			}
			if !reflect.DeepEqual(filesToFormat, tc.expectedFiles) {
				t.Fatalf("Expected to receive paths %v but got %v", tc.expectedFiles, filesToFormat)
			}
		})
	}
}

func formatTempPaths(tempPath string, patterns []string) []string {
	formatted := make([]string, len(patterns))
	for i, pattern := range patterns {
		formatted[i] = fmt.Sprintf("%s/%s", tempPath, pattern)
	}
	return formatted
}
