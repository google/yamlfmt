package yamlfmt

import (
	"os"
	"regexp"

	"github.com/google/yamlfmt/internal/collections"
)

type ContentAnalyzer interface {
	ExcludePathsByContent(paths []string) ([]string, []string, error)
}

type BasicContentAnalyzer struct {
	RegexPatterns []*regexp.Regexp
}

func NewBasicContentAnalyzer(patterns []string) (BasicContentAnalyzer, error) {
	analyzer := BasicContentAnalyzer{RegexPatterns: []*regexp.Regexp{}}
	compileErrs := collections.Errors{}
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			compileErrs = append(compileErrs, err)
			continue
		}
		analyzer.RegexPatterns = append(analyzer.RegexPatterns, re)
	}
	return analyzer, compileErrs.Combine()
}

func (a BasicContentAnalyzer) ExcludePathsByContent(paths []string) ([]string, []string, error) {
	pathsToFormat := collections.SliceToSet(paths)
	pathsExcluded := []string{}
	pathErrs := collections.Errors{}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			pathErrs = append(pathErrs, err)
			continue
		}

		// Search metadata for ignore
		metadata, mdErrs := ReadMetadata(content, path)
		if len(mdErrs) != 0 {
			pathErrs = append(pathErrs, mdErrs...)
		}
		ignoreFound := false
		for md := range metadata {
			if md.Type == MetadataIgnore {
				ignoreFound = true
				break
			}
		}
		if ignoreFound {
			pathsExcluded = append(pathsExcluded, path)
			pathsToFormat.Remove(path)
			continue
		}

		// Check if content matches any regex
		matched := false
		for _, pattern := range a.RegexPatterns {
			if pattern.Match(content) {
				matched = true
			}
		}
		if matched {
			pathsExcluded = append(pathsExcluded, path)
			pathsToFormat.Remove(path)
		}
	}

	return pathsToFormat.ToSlice(), pathsExcluded, pathErrs.Combine()
}
