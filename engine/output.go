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
