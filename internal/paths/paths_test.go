package paths_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/yamlfmt/internal/paths"
)

type dummyPath struct {
	path     string
	filename string
	isDir    bool
}

func (p *dummyPath) Create() error {
	file := p.fullPath()
	var err error
	if p.isDir {
		err = os.Mkdir(file, 0x0700)
	} else {
		err = os.WriteFile(file, []byte{}, 0x0700)
	}
	return err
}

func (p *dummyPath) fullPath() string {
	return filepath.Join(p.path, p.filename)
}

func TestCollectPathsToFormat(t *testing.T) {
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
				file.path = tempPath
				if err := file.Create(); err != nil {
					t.Fatalf("Failed to create file")
				}
			}

			includePaths := formatTempPaths(tempPath, tc.includePatterns)
			if len(includePaths) == 0 {
				includePaths = []string{fmt.Sprintf("%s/**", tempPath)}
			}
			excludePaths := formatTempPaths(tempPath, tc.excludePatterns)
			if len(excludePaths) == 0 {
				excludePaths = []string{}
			}

			paths, err := paths.CollectPathsToFormat(includePaths, excludePaths)
			if err != nil {
				t.Fatalf("CollectPathsToFormat failed: %v", err)
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
