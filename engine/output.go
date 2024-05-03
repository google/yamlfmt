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
	"fmt"

	"github.com/google/yamlfmt"
)

type EngineOutputFormat string

const (
	EngineOutputDefault   EngineOutputFormat = "default"
	EngineOutputSingeLine EngineOutputFormat = "line"
)

func getEngineOutput(t EngineOutputFormat, operation yamlfmt.Operation, files yamlfmt.FileDiffs, quiet bool) (fmt.Stringer, error) {
	switch t {
	case EngineOutputDefault:
		return engineOutput{Operation: operation, Files: files, Quiet: quiet}, nil
	case EngineOutputSingeLine:
		return engineOutputSingleLine{Operation: operation, Files: files, Quiet: quiet}, nil
	}
	return nil, fmt.Errorf("unknown output type: %s", t)
}

type engineOutput struct {
	Operation yamlfmt.Operation
	Files     yamlfmt.FileDiffs
	Quiet     bool
}

func (eo engineOutput) String() string {
	var msg string
	switch eo.Operation {
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
			return "No files will formatted."
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
