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

	"github.com/RageCage64/multilinediff"
	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/paths"
)

type Engine struct {
	Include          []string
	Exclude          []string
	LineSepCharacter string
	Formatter        yamlfmt.Formatter
}

func (e *Engine) FormatAllFiles() error {
	paths, err := paths.CollectPathsToFormat(e.Include, e.Exclude)
	if err != nil {
		return err
	}

	formatErrors := NewFormatFileErrors()
	for _, path := range paths {
		err := e.FormatFile(path)
		if err != nil {
			formatErrors.Add(path, err)
		}
	}

	if formatErrors.Count() > 0 {
		return formatErrors
	}
	return nil
}

func (e *Engine) FormatFile(path string) error {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	formatted, err := e.Formatter.Format(yamlBytes)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, formatted, 0644)
	return err
}

func (e *Engine) LintAllFiles() error {
	paths, err := paths.CollectPathsToFormat(e.Include, e.Exclude)
	if err != nil {
		return err
	}

	lintErrors := NewLintFileErrors()
	for _, path := range paths {
		err := e.LintFile(path)
		if err != nil {
			lintErrors.Add(path, err)
		}
	}

	if lintErrors.Count() > 0 {
		return lintErrors
	}
	return nil
}

func (e *Engine) LintFile(path string) error {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	formatted, err := e.Formatter.Format(yamlBytes)
	if err != nil {
		return err
	}
	diff, diffCount := multilinediff.Diff(string(yamlBytes), string(formatted), e.LineSepCharacter)
	if diffCount > 0 {
		return fmt.Errorf(diff)
	}
	return nil
}

func (e *Engine) DryRunAllFiles() (string, error) {
	paths, err := paths.CollectPathsToFormat(e.Include, e.Exclude)
	if err != nil {
		return "", err
	}

	formatErrors := NewFormatFileErrors()
	dryRunDiffs := NewDryRunDiffs()
	for _, path := range paths {
		diff, diffCount, err := e.DryRunFile(path)
		if err != nil {
			formatErrors.Add(path, err)
		} else if diffCount > 0 {
			dryRunDiffs.Add(path, diff)
		}
	}

	if formatErrors.Count() > 0 {
		return "", formatErrors
	}
	return dryRunDiffs.CombineOutput(), nil
}

func (e *Engine) DryRunFile(path string) (string, int, error) {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return "", 0, err
	}
	formatted, err := e.Formatter.Format(yamlBytes)
	if err != nil {
		return "", 0, err
	}
	diff, diffCount := multilinediff.Diff(string(yamlBytes), string(formatted), e.LineSepCharacter)
	return diff, diffCount, nil
}
