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

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/engine"
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

type Config struct {
	Extensions      []string               `mapstructure:"extensions"`
	Include         []string               `mapstructure:"include"`
	Exclude         []string               `mapstructure:"exclude"`
	Doublestar      bool                   `mapstructure:"doublestar"`
	LineEnding      yamlfmt.LineBreakStyle `mapstructure:"line_ending"`
	FormatterConfig *formatterConfig       `mapstructure:"formatter,omitempty"`
}

type Command struct {
	Operation Operation
	Registry  *yamlfmt.Registry
	Config    *Config
}

func (c *Command) Run() error {
	var formatter yamlfmt.Formatter
	if c.Config.FormatterConfig == nil {
		factory, err := c.Registry.GetDefaultFactory()
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
		if c.Config.FormatterConfig.Type == "" {
			factory, err = c.Registry.GetDefaultFactory()
		} else {
			factory, err = c.Registry.GetFactory(c.Config.FormatterConfig.Type)
		}
		if err != nil {
			return err
		}

		c.Config.FormatterConfig.FormatterSettings["line_ending"] = c.Config.LineEnding
		formatter, err = factory.NewFormatter(c.Config.FormatterConfig.FormatterSettings)
		if err != nil {
			return err
		}
	}

	lineSepChar, err := c.Config.LineEnding.Separator()
	if err != nil {
		return err
	}
	engine := &engine.Engine{
		LineSepCharacter: lineSepChar,
		Formatter:        formatter,
		PathCollector:    c.makePathCollector(),
	}

	switch c.Operation {
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

func (c *Command) makePathCollector() yamlfmt.PathCollector {
	if c.Config.Doublestar {
		return &yamlfmt.DoublestarCollector{
			Include: c.Config.Include,
			Exclude: c.Config.Exclude,
		}
	}
	return &yamlfmt.FilepathCollector{
		Include:    c.Config.Include,
		Exclude:    c.Config.Exclude,
		Extensions: c.Config.Extensions,
	}
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
