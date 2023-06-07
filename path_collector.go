package yamlfmt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/yamlfmt/internal/collections"
)

type PathCollector interface {
	CollectPaths() ([]string, error)
}

type FilepathCollector struct {
	Include    []string
	Exclude    []string
	Extensions []string
}

func (c *FilepathCollector) CollectPaths() ([]string, error) {
	pathsFound := []string{}
	for _, inclPath := range c.Include {
		info, err := os.Stat(inclPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			continue
		}
		if !info.IsDir() {
			pathsFound = append(pathsFound, inclPath)
			continue
		}
		paths, err := c.walkDirectoryForYaml(inclPath)
		if err != nil {
			return nil, err
		}
		pathsFound = append(pathsFound, paths...)
	}

	pathsFoundSet := collections.SliceToSet(pathsFound)
	pathsToFormat := collections.SliceToSet(pathsFound)
	for _, exclPath := range c.Exclude {
		info, err := os.Stat(exclPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			continue
		}

		if info.IsDir() {
			for foundPath := range pathsFoundSet {
				if strings.HasPrefix(foundPath, exclPath) {
					pathsToFormat.Remove(foundPath)
				}
			}
		} else {
			pathsToFormat.Remove(exclPath)
		}
	}

	return pathsToFormat.ToSlice(), nil
}

func (c *FilepathCollector) walkDirectoryForYaml(dir string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		extension := ""
		if strings.Contains(info.Name(), ".") {
			nameParts := strings.Split(info.Name(), ".")
			extension = nameParts[len(nameParts)-1]
		}
		if collections.SliceContains(c.Extensions, extension) {
			paths = append(paths, path)
		}

		return nil
	})
	return paths, err
}

type DoublestarCollector struct {
	Include []string
	Exclude []string
}

func (c *DoublestarCollector) CollectPaths() ([]string, error) {
	includedPaths := []string{}
	for _, pattern := range c.Include {
		globMatches, err := doublestar.FilepathGlob(pattern)
		if err != nil {
			return nil, err
		}
		includedPaths = append(includedPaths, globMatches...)
	}

	pathsToFormatSet := collections.Set[string]{}
	for _, path := range includedPaths {
		if len(c.Exclude) == 0 {
			pathsToFormatSet.Add(path)
			continue
		}
		excluded := false
		for _, pattern := range c.Exclude {
			match, err := doublestar.PathMatch(filepath.Clean(pattern), path)
			if err != nil {
				return nil, err
			}
			if match {
				excluded = true
			}
		}
		if !excluded {
			pathsToFormatSet.Add(path)
		}
	}

	return pathsToFormatSet.ToSlice(), nil
}
