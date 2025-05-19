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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
	"github.com/google/yamlfmt/internal/collections"
)

type GoldenCtx struct {
	Dir    string
	Update bool
}

func (g GoldenCtx) goldenPath(path string) string {
	return filepath.Join(g.Dir, path)
}

func (g GoldenCtx) CompareGoldenFile(path string, gotContent []byte) error {
	// If we are updating, just rewrite the file.
	if g.Update {
		fmt.Println("writing file: ", path)
		return os.WriteFile(g.goldenPath(path), gotContent, os.ModePerm)
	}

	// If we are not updating, check that the content is the same.
	expectedContent, err := os.ReadFile(g.goldenPath(path))
	if err != nil {
		return fmt.Errorf("os.ReadFile(%q): %w", g.goldenPath(path), err)
	}
	// Edge case for empty stdout.
	if gotContent == nil {
		gotContent = []byte{}
	}
	diff := cmp.Diff(string(expectedContent), string(gotContent))
	// If there is no diff between the content, nothing to do in either mode.
	if diff == "" {
		return nil
	}
	return &GoldenDiffError{path: path, diff: diff}
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
	goldenPaths, err := readAllPaths(g.Dir)
	if err != nil {
		return err
	}

	// If the directories differ in paths then the test has failed.
	if !resultPaths.Equals(goldenPaths) {
		return errors.New("the directories were different")
	}

	// Compare each file and gather each error.
	compareErrors := collections.Errors{}
	for path := range resultPaths {
		gotContent, err := os.ReadFile(filepath.Join(resultPath, path))
		if err != nil {
			return fmt.Errorf("os.ReadFile(%q): %w", path, err)
		}
		err = g.CompareGoldenFile(path, gotContent)
		if err != nil {
			return fmt.Errorf("CompareGoldenFile(%q): %w", path, err)
		}
	}
	// If there are no errors this will be nil, otherwise will be a
	// combination of every error that occurred.
	return compareErrors.Combine()
}

func (g GoldenCtx) updateGoldenDirectory(resultPath string) error {
	// Clear the golden directory
	err := os.RemoveAll(g.Dir)
	if err != nil {
		return fmt.Errorf("could not clear golden directory %s: %w", g.Dir, err)
	}
	err = os.Mkdir(g.Dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not recreate golden directory %s: %w", g.Dir, err)
	}

	// Recreate the goldens directory
	paths, err := ReplicateDirectory(resultPath, g.Dir)
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

type GoldenDiffError struct {
	path string
	diff string
}

func (e *GoldenDiffError) Error() string {
	return fmt.Sprintf("golden: %s differed:\n%s", e.path, e.diff)
}
