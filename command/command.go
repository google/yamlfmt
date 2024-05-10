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
	"os"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/engine"

	"github.com/braydonk/yaml"
)

type FormatterConfig struct {
	Type              string         `mapstructure:"type"`
	FormatterSettings map[string]any `mapstructure:",remain"`
}

func (f *FormatterConfig) flatten() map[string]any {
	flat := make(map[string]any, len(f.FormatterSettings)+1)
	if f.Type == "" {
		flat["type"] = yamlfmt.BasicFormatterType
	} else {
		flat["type"] = f.Type
	}

	for k, v := range f.FormatterSettings {
		flat[k] = v
	}

	return flat
}

// NewFormatterConfig returns an empty formatter config with all fields initialized.
func NewFormatterConfig() *FormatterConfig {
	return &FormatterConfig{FormatterSettings: make(map[string]any)}
}

type Config struct {
	Extensions        []string                  `mapstructure:"extensions" yaml:"extensions"`
	Include           []string                  `mapstructure:"include" yaml:"include"`
	Exclude           []string                  `mapstructure:"exclude" yaml:"exclude"`
	RegexExclude      []string                  `mapstructure:"regex_exclude" yaml:"regex_exclude"`
	FormatterConfig   *FormatterConfig          `mapstructure:"formatter,omitempty" yaml:"-"`
	Doublestar        bool                      `mapstructure:"doublestar" yaml:"doublestar"`
	ContinueOnError   bool                      `mapstructure:"continue_on_error" yaml:"continue_on_error"`
	LineEnding        yamlfmt.LineBreakStyle    `mapstructure:"line_ending" yaml:"line_ending"`
	GitignoreExcludes bool                      `mapstructure:"gitignore_excludes" yaml:"gitignore_excludes"`
	GitignorePath     string                    `mapstructure:"gitignore_path" yaml:"gitignore_path"`
	OutputFormat      engine.EngineOutputFormat `mapstructure:"output_format" yaml:"output_format"`
}

func (c *Config) Marshal() ([]byte, error) {
	conf, err := yaml.Marshal(c)
	if err != nil {
		return []byte{}, err
	}
	formatterConf, err := yaml.Marshal(struct {
		Formatter map[string]any `yaml:"formatter"`
	}{Formatter: c.FormatterConfig.flatten()})
	if err != nil {
		return []byte{}, err
	}
	return append(conf, formatterConf...), nil
}

type Command struct {
	Operation yamlfmt.Operation
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
		OutputFormat:     c.Config.OutputFormat,
	}

	collectedPaths, err := c.collectPaths()
	if err != nil {
		return err
	}
	if c.Config.GitignoreExcludes {
		newPaths, err := yamlfmt.ExcludeWithGitignore(c.Config.GitignorePath, collectedPaths)
		if err != nil {
			return err
		}
		collectedPaths = newPaths
	}

	paths, err := c.analyzePaths(collectedPaths)
	if err != nil {
		fmt.Printf("path analysis found the following errors:\n%v", err)
		fmt.Println("Continuing...")
	}

	switch c.Operation {
	case yamlfmt.OperationFormat:
		out, err := eng.Format(paths)
		if out != nil {
			fmt.Print(out)
		}
		if err != nil {
			return err
		}
	case yamlfmt.OperationLint:
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
	case yamlfmt.OperationDry:
		out, err := eng.DryRun(paths)
		if err != nil {
			return err
		}
		if out != nil {
			fmt.Print(out)
		} else if !c.Quiet {
			fmt.Println("No files will be changed.")
		}
	case yamlfmt.OperationStdin:
		stdinYaml, err := readFromStdin()
		if err != nil {
			return err
		}
		out, err := eng.FormatContent(stdinYaml)
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	case yamlfmt.OperationPrintConfig:
		conf, err := c.Config.Marshal()
		if err != nil {
			return err
		}
		formatted, err := formatter.Format(conf)
		if err != nil {
			return err
		}
		fmt.Println(string(formatted))
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
