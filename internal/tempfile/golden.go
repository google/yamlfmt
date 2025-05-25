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

package tempfile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/yamlfmt/internal/collections"
)

type GoldenCtx struct {
	GoldenDir string
	ResultDir string
	Update    bool
}

// Given either a result or golden path, this will strip
// and give you the base.
func (g GoldenCtx) basePath(path string) string {
	if basePath, ok := strings.CutPrefix(path, g.GoldenDir); ok {
		return basePath
	}
	if basePath, ok := strings.CutPrefix(path, g.ResultDir); ok {
		return basePath
	}
	return path
}

func (g GoldenCtx) goldenPath(path string) string {
	return filepath.Join(g.GoldenDir, g.basePath(path))
}

func (g GoldenCtx) CompareGoldenFile(path string, gotContent []byte) error {
	// If we are updating, just rewrite the file.
	if g.Update {
		fmt.Println("writing file: ", path)
		return os.WriteFile(g.goldenPath(path), gotContent, os.ModePerm)
	}

	// If we are not updating, check that the content is the same.
	goldenPath := g.goldenPath(path)
	expectedContent, err := os.ReadFile(goldenPath)
	if err != nil {
		return fmt.Errorf("os.ReadFile(%q): %w", goldenPath, err)
	}
	// Edge case for empty stdout.
	if gotContent == nil {
		gotContent = []byte{}
	}
	diff := cmp.Diff(string(gotContent), string(expectedContent))
	// If there is no diff between the content, nothing to do in either mode.
	if diff == "" {
		return nil
	}
	return &GoldenDiffError{path: g.basePath(path), diff: diff}
}

func (g GoldenCtx) CompareDirectory(resultPath string) error {
	// If in update mode, clobber the whole directory and recreate it with
	// the result of the test.
	if g.Update {
		return g.updateGoldenDirectory(resultPath)
	}

	// Compare the two directories by reading all paths.
	resultPaths, err := readAllPaths(resultPath)
	if err != nil {
		return err
	}
	goldenPaths, err := readAllPaths(g.GoldenDir)
	if err != nil {
		return err
	}

	// If the directories differ in paths then the test has failed.
	if err := directoryFilesEqual(goldenPaths, resultPaths); err != nil {
		return err
	}

	// Compare each file and gather each error.
	compareErrors := collections.Errors{}
	for path := range resultPaths {
		gotContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("os.ReadFile(%q): %w", path, err)
		}
		err = g.CompareGoldenFile(path, gotContent)
		if err != nil {
			return fmt.Errorf("CompareGoldenFile(%q): %w", g.basePath(path), err)
		}
	}
	// If there are no errors this will be nil, otherwise will be a
	// combination of every error that occurred.
	return compareErrors.Combine()
}

func (g GoldenCtx) updateGoldenDirectory(resultPath string) error {
	// Clear the golden directory
	err := os.RemoveAll(g.GoldenDir)
	if err != nil {
		return fmt.Errorf("could not clear golden directory %s: %w", g.GoldenDir, err)
	}
	err = os.Mkdir(g.GoldenDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not recreate golden directory %s: %w", g.GoldenDir, err)
	}

	// Recreate the goldens directory
	paths, err := ReplicateDirectory(resultPath, g.GoldenDir)
	if err != nil {
		return err
	}
	return paths.CreateAll()
}

func readAllPaths(dirPath string) (collections.Set[string], error) {
	paths := collections.Set[string]{}
	allNamesButCurrentDirectory := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		paths.Add(path)
		return nil
	}
	err := filepath.WalkDir(dirPath, allNamesButCurrentDirectory)
	return paths, err
}

func directoryFilesEqual(expectedPaths, actualPaths collections.Set[string]) error {
	expectedFiles := filenamesFromSet(expectedPaths)
	actualFiles := filenamesFromSet(actualPaths)
	if !expectedFiles.Equals(actualFiles) {
		return fmt.Errorf(
			"got different files in generated directory\nexpected: %v\nactual: %v",
			expectedFiles, actualFiles,
		)
	}
	return nil
}

func filenamesFromSet(paths collections.Set[string]) collections.Set[string] {
	files := collections.Set[string]{}
	for path := range paths {
		files.Add(filepath.Base(path))
	}
	return files
}

type GoldenDiffError struct {
	path string
	diff string
}

func (e *GoldenDiffError) Error() string {
	return fmt.Sprintf("golden: %s differed:\n%s", e.path, e.diff)
}
