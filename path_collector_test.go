package yamlfmt_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/yamlfmt"
)

type dummyPath struct {
	basePath string
	filename string
	isDir    bool
}

func (p *dummyPath) Create() error {
	file := p.fullPath()
	var err error
	if p.isDir {
		err = os.Mkdir(file, os.ModePerm)
	} else {
		err = os.WriteFile(file, []byte{}, os.ModePerm)
	}
	return err
}

func (p *dummyPath) fullPath() string {
	return filepath.Join(p.basePath, p.filename)
}

func TestCollectPaths(t *testing.T) {
	testCases := []struct {
		name            string
		files           []dummyPath
		includePatterns []string
		excludePatterns []string
		extensions      []string
		expectedFiles   map[string]struct{}
	}{
		{
			name: "finds direct paths",
			files: []dummyPath{
				{filename: "x.yaml"},
				{filename: "y.yaml"},
				{filename: "z.yml"},
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
			files: []dummyPath{
				{filename: "a", isDir: true},
				{filename: "a/x.yaml"},
				{filename: "a/y.yaml"},
				{filename: "a/z.yml"},
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
			files: []dummyPath{
				{filename: "a", isDir: true},
				{filename: "a/x.yaml"},
				{filename: "a/y.yaml"},
				{filename: "a/z.yml"},
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
			files: []dummyPath{
				{filename: "a", isDir: true},
				{filename: "a/b", isDir: true},
				{filename: "x.yml"},
				{filename: "y.yml"},
				{filename: "z.yaml"},
				{filename: "a/x.yaml"},
				{filename: "a/b/x.yaml"},
				{filename: "a/b/y.yml"},
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
			files: []dummyPath{
				{filename: "a", isDir: true},
				{filename: "a/b", isDir: true},
				{filename: "x.yml"},
				{filename: "y.yml"},
				{filename: "z.yaml"},
				{filename: "a/x.yaml"},
				{filename: "a/b/x.yaml"},
				{filename: "a/b/y.yml"},
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
			files: []dummyPath{
				{filename: "a", isDir: true},
				{filename: "a/b", isDir: true},
				{filename: "x.yml"},
				{filename: "y.yml"},
				{filename: "z.yaml"},
				{filename: "a/x.yaml"},
				{filename: "a/b/x.yaml"},
				{filename: "a/b/y.yml"},
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
				file.basePath = tempPath
				if err := file.Create(); err != nil {
					t.Fatalf("Failed to create file %s: %v", file.basePath, err)
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
		files           []dummyPath
		includePatterns []string
		excludePatterns []string
		expectedFiles   map[string]struct{}
	}{
		{
			name: "no excludes",
			files: []dummyPath{
				{filename: "x.yaml"},
				{filename: "y.yaml"},
				{filename: "z.yaml"},
			},
			expectedFiles: map[string]struct{}{
				"x.yaml": {},
				"y.yaml": {},
				"z.yaml": {},
			},
		},
		{
			name: "does not include directories",
			files: []dummyPath{
				{filename: "x.yaml"},
				{filename: "y", isDir: true},
			},
			expectedFiles: map[string]struct{}{
				"x.yaml": {},
			},
		},
		{
			name: "only include what is asked",
			files: []dummyPath{
				{filename: "x.yaml"},
				{filename: "y.yaml"},
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
			files: []dummyPath{
				{filename: "x.yaml"},
				{filename: "y.yaml"},
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
				file.basePath = tempPath
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
