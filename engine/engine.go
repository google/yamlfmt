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

type Engine struct {
	Include   []string
	Exclude   []string
	Formatter yamlfmt.Formatter
}

func (e *Engine) FormatAllFiles() error {
	paths, err := CollectPathsToFormat(e.Include, e.Exclude)
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

func (f *Engine) LintAllFiles() error {
	paths, err := CollectPathsToFormat(f.Include, f.Exclude)
	if err != nil {
		return err
	}

	lintErrors := NewLintFileErrors()
	for _, path := range paths {
		err := f.LintFile(path)
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
	diff := MultilineStringDiff(string(yamlBytes), string(formatted))
	if diff != "" {
		return fmt.Errorf(diff)
	}
	return nil
}

func (f *Engine) DryRunAllFiles() (string, error) {
	paths, err := CollectPathsToFormat(f.Include, f.Exclude)
	if err != nil {
		return "", err
	}

	formatErrors := NewFormatFileErrors()
	dryRunDiffs := NewDryRunDiffs()
	for _, path := range paths {
		diff, err := f.DryRunFile(path)
		if err != nil {
			formatErrors.Add(path, err)
		} else if diff != "" {
			dryRunDiffs.Add(path, diff)
		}
	}

	if formatErrors.Count() > 0 {
		return "", formatErrors
	}
	return dryRunDiffs.CombineOutput(), nil
}

func (e *Engine) DryRunFile(path string) (string, error) {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	formatted, err := e.Formatter.Format(yamlBytes)
	if err != nil {
		return "", err
	}
	diff := MultilineStringDiff(string(yamlBytes), string(formatted))
	return diff, nil
}
