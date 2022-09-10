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
	"github.com/google/yamlfmt/internal/hotfix"
	"gopkg.in/yaml.v3"
)

const BasicFormatterType string = "basic"

type BasicFormatter struct {
	Config   *Config
	Features yamlfmt.FeatureList
}

func (f *BasicFormatter) ConfigureFeaturesFromConfig() {
	if f.Config.EmojiSupport {
		f.Features = append(f.Features, featEmojiSupport)
	}
	if f.Config.IncludeDocumentStart {
		f.Features = append(f.Features, featIncludeDocumentStart)
	}
	if f.Config.LineEnding == yamlfmt.LineBreakStyleCRLF {
		f.Features = append(f.Features, featCRLFSupport)
	}
	if f.Config.RetainLineBreaks {
		linebreakStr := "\n"
		if f.Config.LineEnding == yamlfmt.LineBreakStyleCRLF {
			linebreakStr = "\r\n"
		}
		featLineBreak := hotfix.MakeFeatureRetainLineBreak(linebreakStr, f.Config.Indent)
		f.Features = append(f.Features, featLineBreak)
	}
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
