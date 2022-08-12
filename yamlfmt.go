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

package yamlfmt

import "fmt"

type BaseConfig struct {
	Include []string `mapstructure:"include"`
	Exclude []string `mapstructure:"exclude"`
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Include: []string{"**/*.yaml", "**/*.yml"},
		Exclude: []string{},
	}
}

type Formatter interface {
	Type() string
	FormatAllFiles() error
	FormatFile(path string) error
	LintAllFiles() error
	LintFile(path string) error
	Format(yamlContent []byte) ([]byte, error)
}

type Factory interface {
	NewWithConfig(config map[string]interface{}) (Formatter, error)
	NewDefault() Formatter
}

type FormatterRegistry struct {
	registry map[string]Factory
}

func NewFormatterRegistry() *FormatterRegistry {
	return &FormatterRegistry{registry: map[string]Factory{}}
}

func (r *FormatterRegistry) Add(fType string, f Factory) {
	r.registry[fType] = f
}

func (r *FormatterRegistry) GetFactory(fType string) (Factory, error) {
	factory, ok := r.registry[fType]
	if !ok {
		return nil, fmt.Errorf("no formatter registered with type \"%s\"", fType)
	}
	return factory, nil
}
