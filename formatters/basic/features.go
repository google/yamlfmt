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
	"github.com/google/yamlfmt/internal/hotfix"
)

var (
	featIncludeDocumentStart = yamlfmt.Feature{
		Name: "Include Document Start",
		AfterAction: func(content []byte) ([]byte, error) {
			documentStart := "---\n"
			return append([]byte(documentStart), content...), nil
		},
	}
	featEmojiSupport = yamlfmt.Feature{
		Name:        "Emoji Support",
		AfterAction: hotfix.ParseUnicodePoints,
	}
	featCRLFSupport = yamlfmt.Feature{
		Name:         "CRLF Support",
		BeforeAction: hotfix.StripCRBytes,
		AfterAction:  hotfix.WriteCRLFBytes,
	}
)

func ConfigureFeaturesFromConfig(config *Config) yamlfmt.FeatureList {
	features := []yamlfmt.Feature{}
	if config.EmojiSupport {
		features = append(features, featEmojiSupport)
	}
	if config.IncludeDocumentStart {
		features = append(features, featIncludeDocumentStart)
	}
	if config.LineEnding == yamlfmt.LineBreakStyleCRLF {
		features = append(features, featCRLFSupport)
	}
	if config.RetainLineBreaks {
		lineSep, err := config.LineEnding.Separator()
		if err != nil {
			lineSep = "\n"
		}
		featLineBreak := hotfix.MakeFeatureRetainLineBreak(lineSep)
		features = append(features, featLineBreak)
	}
	return features
}
