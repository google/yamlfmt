package yamlfmt

import (
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
		for _, pattern := range exclude {
			match, err := doublestar.Match(pattern, path)
			if err != nil {
				return nil, err
			}
			if !match {
				pathsToFormatSet[path] = struct{}{}
			}
		}
	}
	pathsToFormat := []string{}
	for path := range pathsToFormatSet {
		pathsToFormat = append(pathsToFormat, path)
	}
	return pathsToFormat, nil
}
