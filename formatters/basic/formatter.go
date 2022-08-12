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

package basic

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/errors"
)

const BasicFormatterType string = "basic"

type Formatter struct {
	config Config
}

func (f *Formatter) Type() string {
	return BasicFormatterType
}

func (f *Formatter) FormatAllFiles() error {
	paths, err := yamlfmt.CollectPathsToFormat(f.config.Include, f.config.Exclude)
	if err != nil {
		return err
	}

	formatErrors := errors.NewFormatFileErrors()
	for _, path := range paths {
		err := f.FormatFile(path)
		if err != nil {
			formatErrors.Add(path, err)
		}
	}

	if formatErrors.Count() > 0 {
		return formatErrors
	}
	return nil
}

func (f *Formatter) FormatFile(path string) error {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	formatted, err := f.Format(yamlBytes)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, formatted, 0644)
	return err
}

func (f *Formatter) LintAllFiles() error {
	paths, err := yamlfmt.CollectPathsToFormat(f.config.Include, f.config.Exclude)
	if err != nil {
		return err
	}

	lintErrors := errors.NewLintFileErrors()
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

func (f *Formatter) LintFile(path string) error {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	formatted, err := f.Format(yamlBytes)
	if err != nil {
		return err
	}
	diff := cmp.Diff(
		string(yamlBytes), string(formatted),
		// Diff each line separately for readability.
		cmpopts.AcyclicTransformer("multiline", func(s string) []string {
			return strings.Split(s, "\n")
		}),
	)
	if diff != "" {
		return fmt.Errorf(diff)
	}
	return nil
}

func (f *Formatter) Format(yamlContent []byte) ([]byte, error) {
	var unmarshalled map[string]interface{}
	err := yaml.Unmarshal(yamlContent, &unmarshalled)
	if err != nil {
		return nil, err
	}
	marshalled, err := yaml.Marshal(unmarshalled)
	if err != nil {
		return nil, err
	}
	return marshalled, nil
}
