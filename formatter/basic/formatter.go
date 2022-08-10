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
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Formatter struct {
	config *Config
}

func (f *Formatter) Type() string {
	return "basic"
}

func (f *Formatter) FormatAllFiles() error {
	paths, err := f.collectAllPaths()
	if err != nil {
		return err
	}

	errors := &formatFileErrors{}
	for _, path := range paths {
		err := f.Format(path)
		if err != nil {
			errors.formatErrors[path] = err
		}
	}

	if len(errors.formatErrors) > 0 {
		return errors
	}
	return nil
}

func (f *Formatter) Format(path string) error {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var unmarshalled map[string]interface{}
	err = yaml.UnmarshalStrict(yamlBytes, &unmarshalled)
	if err != nil {
		return err
	}
	marshalled, err := yaml.Marshal(unmarshalled)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, marshalled, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (f *Formatter) collectAllPaths() ([]string, error) {
	includedPaths := []string{}
	for _, pattern := range f.config.Include {
		globMatches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		includedPaths = append(includedPaths, globMatches...)
	}

	pathsToFormat := []string{}
	for _, path := range includedPaths {
		for _, pattern := range f.config.Exclude {
			match, err := filepath.Match(pattern, path)
			if err != nil {
				return nil, err
			}
			if !match {
				pathsToFormat = append(pathsToFormat, path)
			}
		}
	}
	return pathsToFormat, nil
}
