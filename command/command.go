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
	"errors"
	"fmt"
	"io"
	"log"
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

type FormatterConfig struct {
	Type              string         `mapstructure:"type"`
	FormatterSettings map[string]any `mapstructure:",remain"`
}

// NewFormatterConfig returns an empty formatter config with all fields initialized.
func NewFormatterConfig() FormatterConfig {
	return FormatterConfig{FormatterSettings: make(map[string]any)}
}

type Config struct {
	Extensions      []string               `mapstructure:"extensions"`
	Include         []string               `mapstructure:"include"`
	Exclude         []string               `mapstructure:"exclude"`
	RegexExclude    []string               `mapstructure:"regex_exclude"`
	Doublestar      bool                   `mapstructure:"doublestar"`
	ContinueOnError bool                   `mapstructure:"continue_on_error"`
	LineEnding      yamlfmt.LineBreakStyle `mapstructure:"line_ending"`
	FormatterConfig *FormatterConfig       `mapstructure:"formatter,omitempty"`
}

// NewConfig returns an empty config with all fields initialized.
func NewConfig() Config {
	formatterConfig := NewFormatterConfig()
	return Config{FormatterConfig: &formatterConfig}
}

type Command struct {
	Operation Operation
	Registry  *yamlfmt.Registry
	Config    *Config
	Quiet     bool
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

	eng := &engine.ConsecutiveEngine{
		LineSepCharacter: lineSepChar,
		Formatter:        formatter,
		Quiet:            c.Quiet,
		ContinueOnError:  c.Config.ContinueOnError,
	}

	collectedPaths, err := c.collectPaths()
	if err != nil {
		return err
	}
	paths, err := c.analyzePaths(collectedPaths)
	if err != nil {
		log.Printf("path analysis found the following errors:\n%v", err)
		log.Println("Continuing...")
	}

	switch c.Operation {
	case OperationFormat:
		err := eng.Format(paths)
		if err != nil {
			return err
		}
	case OperationLint:
		out, err := eng.Lint(paths)
		if err != nil {
			return err
		}
		if out != nil {
			// This will be picked up by log.Fatal in main() and
			// cause an exit code of 1, which is a critical
			// component of the lint functionality.
			return errors.New(out.String())
		}
	case OperationDry:
		out, err := eng.DryRun(paths)
		if err != nil {
			return err
		}
		if out.Message == "" && out.Files.ChangedCount() == 0 {
			log.Print("No files will be changed.")
		} else {
			log.Print(out)
		}
	case OperationStdin:
		stdinYaml, err := readFromStdin()
		if err != nil {
			return err
		}
		out, err := eng.FormatContent(stdinYaml)
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	}

	return nil
}

func (c *Command) collectPaths() ([]string, error) {
	collector := c.makePathCollector()
	return collector.CollectPaths()
}

func (c *Command) analyzePaths(paths []string) ([]string, error) {
	analyzer, err := c.makeAnalyzer()
	if err != nil {
		return nil, err
	}
	includePaths, _, err := analyzer.ExcludePathsByContent(paths)
	return includePaths, err
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

func (c *Command) makeAnalyzer() (yamlfmt.ContentAnalyzer, error) {
	return yamlfmt.NewBasicContentAnalyzer(c.Config.RegexExclude)
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
