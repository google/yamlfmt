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
	"github.com/google/yamlfmt"
	"github.com/mitchellh/mapstructure"
)

type BasicFormatterFactory struct{}

func (f *BasicFormatterFactory) Type() string {
	return BasicFormatterType
}

func (f *BasicFormatterFactory) NewDefault() yamlfmt.Formatter {
	return &BasicFormatter{Config: DefaultConfig()}
}

func (f *BasicFormatterFactory) NewWithConfig(configData map[string]interface{}) (yamlfmt.Formatter, error) {
	var c Config
	err := mapstructure.Decode(configData, &c)
	if err != nil {
		return nil, err
	}
	return &BasicFormatter{Config: &c}, nil
}
