// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package paths

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func CollectPathsToFormat(include, exclude []string) ([]string, error) {
	includedPaths := []string{}
	for _, pattern := range include {
		globMatches, err := doublestar.FilepathGlob(pattern)
		if err != nil {
			return nil, err
		}
		includedPaths = append(includedPaths, globMatches...)
	}

	pathsToFormatSet := map[string]struct{}{}
	for _, path := range includedPaths {
		if len(exclude) == 0 {
			pathsToFormatSet[path] = struct{}{}
			continue
		}
		excluded := false
		for _, pattern := range exclude {
			match, err := doublestar.Match(filepath.Clean(pattern), path)
			if err != nil {
				return nil, err
			}
			if match {
				excluded = true
			}
		}
		if !excluded {
			pathsToFormatSet[path] = struct{}{}
		}
	}
	pathsToFormat := []string{}
	for path := range pathsToFormatSet {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		pathsToFormat = append(pathsToFormat, path)
	}
	return pathsToFormat, nil
}
