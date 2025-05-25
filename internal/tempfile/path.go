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
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/yamlfmt/internal/collections"
)

type Path struct {
	// BasePath is the new location for the file represented by this Path.
	BasePath string
	// FilePath is the destination path within the BasePath.
	FilePath string
	IsDir    bool
	Content  []byte
}

func (p *Path) Create() error {
	file := p.fullPath()
	var err error
	if p.IsDir {
		err = os.Mkdir(file, os.ModePerm)
	} else {
		err = os.WriteFile(file, p.Content, os.ModePerm)
	}
	return err
}

func (p *Path) fullPath() string {
	return filepath.Join(p.BasePath, p.FilePath)
}

type Paths []Path

func (ps Paths) CreateAll() error {
	errs := collections.Errors{}
	for _, path := range ps {
		errs = append(errs, path.Create())
	}
	return errs.Combine()
}

func ReplicateDirectory(dir string, newBase string) (Paths, error) {
	paths := Paths{}
	walkAllButCurrentDirectory := func(path string, info fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		// Skip the current directory (basically the . directory)
		if path == dir {
			return nil
		}
		content := []byte{}

		if !info.IsDir() {
			readContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content = readContent
		}

		paths = append(paths, Path{
			BasePath: newBase,
			FilePath: strings.TrimPrefix(path, dir),
			IsDir:    info.IsDir(),
			Content:  content,
		})
		return nil
	}
	err := filepath.WalkDir(dir, walkAllButCurrentDirectory)
	return paths, err
}
