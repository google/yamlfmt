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

package engine

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/gitlab"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type EngineOutputFormat string

const (
	EngineOutputDefault   EngineOutputFormat = "default"
	EngineOutputSingeLine EngineOutputFormat = "line"
	EngineOutputSlim      EngineOutputFormat = "slim"
	EngineOutputGitlab    EngineOutputFormat = "gitlab"
)

func getEngineOutput(t EngineOutputFormat, operation yamlfmt.Operation, files yamlfmt.FileDiffs, quiet bool, verbose bool) (fmt.Stringer, error) {
	switch t {
	case EngineOutputDefault:
		return engineOutput{Operation: operation, Files: files, Quiet: quiet, Verbose: verbose}, nil
	case EngineOutputSingeLine:
		return engineOutputSingleLine{Operation: operation, Files: files, Quiet: quiet}, nil
	case EngineOutputSlim:
		return engineOutputSlim{Operation: operation, Files: files, Quiet: quiet, Verbose: verbose}, nil
	case EngineOutputGitlab:
		return engineOutputGitlab{Operation: operation, Files: files, Compact: quiet}, nil

	}
	return nil, fmt.Errorf("unknown output type: %s", t)
}

type engineOutput struct {
	Operation yamlfmt.Operation
	Files     yamlfmt.FileDiffs
	Quiet     bool
	Verbose   bool
}

func (eo engineOutput) String() string {
	var msg string
	switch eo.Operation {
	case yamlfmt.OperationFormat:
		// Formatting only produces output in verbose mode.
		if !eo.Verbose {
			return ""
		}

		msg = "The following files were modified:\n"
		msg += eo.Files.StrOutputQuiet()
		return msg
	case yamlfmt.OperationLint:
		msg = "The following formatting differences were found:"
		if eo.Quiet {
			msg = "The following files had formatting differences:"
		}
	case yamlfmt.OperationDry:
		if len(eo.Files) > 0 {
			if eo.Quiet {
				msg = "The following files would be formatted:"
			}
		} else {
			return "No files will be formatted."
		}
	}
	var result string
	if msg != "" {
		result += fmt.Sprintf("%s\n\n", msg)
	}
	if eo.Quiet {
		result += eo.Files.StrOutputQuiet()
	} else {
		result += fmt.Sprintf("%s\n", eo.Files.StrOutput())
	}
	return result
}

type engineOutputSingleLine struct {
	Operation yamlfmt.Operation
	Files     yamlfmt.FileDiffs
	Quiet     bool
}

func (eosl engineOutputSingleLine) String() string {
	var msg string
	for _, fileDiff := range eosl.Files {
		msg += fmt.Sprintf("%s: formatting difference found\n", fileDiff.Path)
	}
	return msg
}

type engineOutputSlim struct {
	Operation yamlfmt.Operation
	Files     yamlfmt.FileDiffs
	Quiet     bool
	Verbose   bool
}

func (eos engineOutputSlim) String() string {
	var msg string
	switch eos.Operation {
	case yamlfmt.OperationFormat:
		// Formatting only produces output in verbose mode.
		if !eos.Verbose {
			return ""
		}
		msg = "The following files were modified:"
	case yamlfmt.OperationLint:
		msg = "The following formatting differences were found:"
		if eos.Quiet {
			msg = "The following files had formatting differences:"
		}
	case yamlfmt.OperationDry:
		if eos.Files.ChangedCount() == 0 {
			return "No files will be formatted."
		}
		msg = "The following files would be formatted:"
	}

	var b strings.Builder
	if msg != "" {
		b.WriteString(msg + "\n\n")
	}

	if eos.Quiet {
		b.WriteString(eos.Files.StrOutputQuiet())
		return b.String()
	}

	for _, fileDiff := range eos.changedFileDiffs() {
		output := slimFileDiff(fileDiff)
		b.WriteString(output)
	}
	return b.String()
}

func (eos engineOutputSlim) changedFileDiffs() []*yamlfmt.FileDiff {
	changed := make([]*yamlfmt.FileDiff, 0, eos.Files.ChangedCount())
	for _, path := range eos.Files.SortedPaths() {
		fileDiff := eos.Files[path]
		if !fileDiff.Diff.Changed() {
			continue
		}
		changed = append(changed, fileDiff)
	}
	return changed
}

func slimFileDiff(fileDiff *yamlfmt.FileDiff) string {
	original := fileDiff.Diff.GetOriginal()
	formatted := fileDiff.Diff.GetFormatted()
	edits := myers.ComputeEdits(span.URIFromPath(fileDiff.Path), original, formatted)
	return fmt.Sprint(gotextdiff.ToUnified(fileDiff.Path, fileDiff.Path, original, edits))
}

type engineOutputGitlab struct {
	Operation yamlfmt.Operation
	Files     yamlfmt.FileDiffs
	Compact   bool
}

func (eo engineOutputGitlab) String() string {
	var findings []gitlab.CodeQuality

	for _, file := range eo.Files {
		if cq, ok := gitlab.NewCodeQuality(*file); ok {
			findings = append(findings, cq)
		}
	}

	if len(findings) == 0 {
		return ""
	}

	sort.Sort(byPath(findings))

	var b strings.Builder
	enc := json.NewEncoder(&b)

	if !eo.Compact {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(findings); err != nil {
		panic(err)
	}
	return b.String()
}

// byPath is used to sort by Location.Path.
type byPath []gitlab.CodeQuality

func (b byPath) Len() int           { return len(b) }
func (b byPath) Less(i, j int) bool { return b[i].Location.Path < b[j].Location.Path }
func (b byPath) Swap(i, j int) {
	b[i].Location.Path, b[j].Location.Path = b[j].Location.Path, b[i].Location.Path
}
