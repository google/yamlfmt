package engine

import (
	"fmt"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func MultilineStringDiff(a, b string) string {
	return cmp.Diff(
		a, b,
		cmpopts.AcyclicTransformer("multiline", func(s string) []string {
			return strings.Split(s, "\n")
		}),
	)
}

type DryRunDiffs struct {
	diffs map[string]string
}

func NewDryRunDiffs() *DryRunDiffs {
	return &DryRunDiffs{diffs: map[string]string{}}
}

func (d *DryRunDiffs) Add(path, diff string) {
	d.diffs[path] = diff
}

func (d *DryRunDiffs) CombineOutput() string {
	if len(d.diffs) == 0 {
		return "dry run produced no output"
	}

	s := ""
	for path, diff := range d.diffs {
		s += fmt.Sprintf("%s:\n%s\n", path, diff)
	}
	return s
}
