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
	"runtime"

	"github.com/google/yamlfmt"
)

type Config struct {
	Indent               int                    `mapstructure:"indent"`
	IncludeDocumentStart bool                   `mapstructure:"include_document_start"`
	LineEnding           yamlfmt.LineBreakStyle `mapstructure:"line_ending"`
	LineLength           int                    `mapstructure:"max_line_length"`
	RetainLineBreaks     bool                   `mapstructure:"retain_line_breaks"`
	DisallowAnchors      bool                   `mapstructure:"disallow_anchors"`
	ScanFoldedAsLiteral  bool                   `mapstructure:"scan_folded_as_literal"`
	IndentlessArrays     bool                   `mapstructure:"indentless_arrays"`
	DropMergeTag         bool                   `mapstructure:"drop_merge_tag"`
	PadLineComments      int                    `mapstructure:"pad_line_comments"`
}

func DefaultConfig() *Config {
	lineBreakStyle := yamlfmt.LineBreakStyleLF
	if runtime.GOOS == "windows" {
		lineBreakStyle = yamlfmt.LineBreakStyleCRLF
	}
	return &Config{
		Indent:          2,
		LineEnding:      lineBreakStyle,
		PadLineComments: 1,
	}
}
