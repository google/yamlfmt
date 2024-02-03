package yamlfmt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/yamlfmt/internal/collections"
	"github.com/google/yamlfmt/internal/logger"
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
	logger.Debug(logger.DebugCodePaths, "using file path matching. include patterns: %s", c.Include)
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
	logger.Debug(logger.DebugCodePaths, "found paths: %s", pathsFound)

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
			logger.Debug(logger.DebugCodePaths, "for exclude dir: %s", exclPath)
			for foundPath := range pathsFoundSet {
				if strings.HasPrefix(foundPath, exclPath) {
					logger.Debug(logger.DebugCodePaths, "excluding %s", foundPath)
					pathsToFormat.Remove(foundPath)
				}
			}
		} else {
			logger.Debug(logger.DebugCodePaths, "for exclude file: %s", exclPath)
			removed := pathsToFormat.Remove(exclPath)
			if removed {
				logger.Debug(logger.DebugCodePaths, "found in paths, excluding")
			}
		}
	}

	pathsToFormatSlice := pathsToFormat.ToSlice()
	logger.Debug(logger.DebugCodePaths, "paths to format: %s", pathsToFormatSlice)
	return pathsToFormatSlice, nil
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
	logger.Debug(logger.DebugCodePaths, "using doublestar path matching. include patterns: %s", c.Include)
	includedPaths := []string{}
	for _, pattern := range c.Include {
		logger.Debug(logger.DebugCodePaths, "trying pattern: %s", pattern)
		globMatches, err := doublestar.FilepathGlob(pattern)
		if err != nil {
			return nil, err
		}
		logger.Debug(logger.DebugCodePaths, "pattern %s matches: %s", pattern, globMatches)
		includedPaths = append(includedPaths, globMatches...)
	}

	pathsToFormatSet := collections.Set[string]{}
	for _, path := range includedPaths {
		if len(c.Exclude) == 0 {
			pathsToFormatSet.Add(path)
			continue
		}
		excluded := false
		logger.Debug(logger.DebugCodePaths, "calculating excludes for %s", path)
		for _, pattern := range c.Exclude {
			match, err := doublestar.PathMatch(filepath.Clean(pattern), path)
			if err != nil {
				return nil, err
			}
			if match {
				logger.Debug(logger.DebugCodePaths, "pattern %s matched, excluding", pattern)
				excluded = true
				break
			}
			logger.Debug(logger.DebugCodePaths, "pattern %s did not match path", pattern)
		}
		if !excluded {
			logger.Debug(logger.DebugCodePaths, "path %s included", path)
			pathsToFormatSet.Add(path)
		}
	}

	pathsToFormat := pathsToFormatSet.ToSlice()
	logger.Debug(logger.DebugCodePaths, "paths to format: %s", pathsToFormat)
	return pathsToFormat, nil
}
