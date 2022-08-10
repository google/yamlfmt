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
	"path"

	"github.com/google/yamlfmt/formatter"
	"gopkg.in/yaml.v2"
)

const defaultConfigName = "yamlfmt.yaml"

type Config struct {
	formatter.BaseConfig
}

func NewDefaultConfig() *Config {
	return &Config{
		Include: []string{"*"},
	}
}

func NewConfigFromFile() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configPath := path.Join(wd, defaultConfigName)
	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}
	yamlBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.UnmarshalStrict(yamlBytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
