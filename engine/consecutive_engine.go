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

package engine

import (
	"fmt"
	"os"

	"github.com/google/yamlfmt"
)

// Engine that will process each file one by one consecutively.
type ConsecutiveEngine struct {
	LineSepCharacter string
	Formatter        yamlfmt.Formatter
	Quiet            bool
	ContinueOnError  bool
	OutputFormat     EngineOutputFormat
}

func (e *ConsecutiveEngine) FormatContent(content []byte) ([]byte, error) {
	return e.Formatter.Format(content)
}

func (e *ConsecutiveEngine) Format(paths []string) (fmt.Stringer, error) {
	formatDiffs, formatErrs := e.formatAll(paths)
	if len(formatErrs) > 0 {
		if e.ContinueOnError {
			fmt.Print(formatErrs)
			fmt.Println("Continuing...")
		} else {
			return nil, formatErrs
		}
	}
	return nil, formatDiffs.ApplyAll()
}

func (e *ConsecutiveEngine) Lint(paths []string) (fmt.Stringer, error) {
	formatDiffs, formatErrs := e.formatAll(paths)
	if len(formatErrs) > 0 {
		return nil, formatErrs
	}
	if formatDiffs.ChangedCount() == 0 {
		return nil, nil
	}
	return getEngineOutput(e.OutputFormat, yamlfmt.OperationLint, formatDiffs, e.Quiet)
}

func (e *ConsecutiveEngine) DryRun(paths []string) (fmt.Stringer, error) {
	formatDiffs, formatErrs := e.formatAll(paths)
	if len(formatErrs) > 0 {
		return nil, formatErrs
	}
	if formatDiffs.ChangedCount() == 0 {
		return nil, nil
	}
	return getEngineOutput(e.OutputFormat, yamlfmt.OperationDry, formatDiffs, e.Quiet)
}

func (e *ConsecutiveEngine) formatAll(paths []string) (yamlfmt.FileDiffs, FormatErrors) {
	formatDiffs := yamlfmt.FileDiffs{}
	formatErrs := FormatErrors{}
	for _, path := range paths {
		fileDiff, err := e.formatFileContent(path)
		if err != nil {
			formatErrs = append(formatErrs, wrapFormatError(path, err))
			continue
		}
		formatDiffs.Add(fileDiff)
	}
	return formatDiffs, formatErrs
}

func (e *ConsecutiveEngine) formatFileContent(path string) (*yamlfmt.FileDiff, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	formatted, err := e.FormatContent(content)
	if err != nil {
		return nil, err
	}
	return &yamlfmt.FileDiff{
		Path: path,
		Diff: &yamlfmt.FormatDiff{
			Original:  string(content),
			Formatted: string(formatted),
			LineSep:   e.LineSepCharacter,
		},
	}, nil
}
