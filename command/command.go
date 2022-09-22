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

package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/engine"
	"github.com/mitchellh/mapstructure"
)

type Operation int

const (
	OperationFormat Operation = iota
	OperationLint
	OperationDry
	OperationStdin
)

type formatterConfig struct {
	Type              string                 `mapstructure:"type"`
	FormatterSettings map[string]interface{} `mapstructure:",remain"`
}

type commandConfig struct {
	Include         []string               `mapstructure:"include"`
	Exclude         []string               `mapstructure:"exclude"`
	LineEnding      yamlfmt.LineBreakStyle `mapstructure:"line_ending"`
	FormatterConfig *formatterConfig       `mapstructure:"formatter,omitempty"`
}

func RunCommand(
	operation Operation,
	registry *yamlfmt.Registry,
	configData map[string]interface{},
) error {
	var config commandConfig
	err := mapstructure.Decode(configData, &config)
	if err != nil {
		return err
	}
	if len(config.Include) == 0 {
		config.Include = []string{"**/*.{yaml,yml}"}
	}
	if config.LineEnding == "" {
		config.LineEnding = yamlfmt.LineBreakStyleLF
		if runtime.GOOS == "windows" {
			config.LineEnding = yamlfmt.LineBreakStyleCRLF
		}
	}

	var formatter yamlfmt.Formatter
	if config.FormatterConfig == nil {
		factory, err := registry.GetDefaultFactory()
		if err != nil {
			return err
		}
		formatter, err = factory.NewFormatter(nil)
		if err != nil {
			return err
		}
	} else {
		var (
			factory yamlfmt.Factory
			err     error
		)
		if config.FormatterConfig.Type == "" {
			factory, err = registry.GetDefaultFactory()
		} else {
			factory, err = registry.GetFactory(config.FormatterConfig.Type)
		}
		if err != nil {
			return err
		}

		config.FormatterConfig.FormatterSettings["line_ending"] = config.LineEnding
		formatter, err = factory.NewFormatter(config.FormatterConfig.FormatterSettings)
		if err != nil {
			return err
		}
	}

	lineSepChar, err := config.LineEnding.Separator()
	if err != nil {
		return err
	}
	engine := &engine.Engine{
		Include:          config.Include,
		Exclude:          config.Exclude,
		LineSepCharacter: lineSepChar,
		Formatter:        formatter,
	}

	switch operation {
	case OperationFormat:
		err := engine.FormatAllFiles()
		if err != nil {
			return err
		}
	case OperationLint:
		err := engine.LintAllFiles()
		if err != nil {
			return err
		}
	case OperationDry:
		out, err := engine.DryRunAllFiles()
		if err != nil {
			return err
		}
		fmt.Println(out)
	case OperationStdin:
		stdinYaml, err := readFromStdin()
		if err != nil {
			return err
		}
		out, err := engine.Formatter.Format(stdinYaml)
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	}

	return nil
}

func readFromStdin() ([]byte, error) {
	stdin := bufio.NewReader(os.Stdin)
	data := []byte{}
	for {
		b, err := stdin.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		data = append(data, b)
	}
	return data, nil
}
