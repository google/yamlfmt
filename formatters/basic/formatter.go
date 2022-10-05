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
	"bytes"
	"errors"
	"io"

	"github.com/google/yamlfmt"
	"gopkg.in/yaml.v3"
)

const BasicFormatterType string = "basic"

type BasicFormatter struct {
	Config   *Config
	Features yamlfmt.FeatureList
}

// yamlfmt.Formatter interface

func (f *BasicFormatter) Type() string {
	return BasicFormatterType
}

func (f *BasicFormatter) Format(input []byte) ([]byte, error) {
	// Run all featurres with BeforeActions
	yamlContent, err := f.Features.ApplyFeatures(input, yamlfmt.FeatureApplyBefore)
	if err != nil {
		return nil, err
	}

	// Format the yaml content
	reader := bytes.NewReader(yamlContent)
	decoder := yaml.NewDecoder(reader)
	documents := []yaml.Node{}
	for {
		var docNode yaml.Node
		err := decoder.Decode(&docNode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		documents = append(documents, docNode)
	}

	// Run all features with DuringActions.
	for _, d := range documents {
		if err := f.Features.ApplyYAMLFeatures(d); err != nil {
			return nil, err
		}
	}

	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(f.Config.Indent)
	for _, doc := range documents {
		err := e.Encode(&doc)
		if err != nil {
			return nil, err
		}
	}

	// Run all features with AfterActions
	resultYaml, err := f.Features.ApplyFeatures(b.Bytes(), yamlfmt.FeatureApplyAfter)
	if err != nil {
		return nil, err
	}

	return resultYaml, nil
}
