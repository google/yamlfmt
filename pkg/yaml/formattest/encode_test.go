// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package formattest

import (
	"testing"

	"github.com/google/yamlfmt/pkg/yaml"
)

func TestExplicitDocumentStart(t *testing.T) {
	formatTestCase{
		name:             "explicit document start",
		folder:           "document_start",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetExplicitDocumentStart(true)
		},
	}.Run(t)
}

func TestIndentless(t *testing.T) {
	formatTestCase{
		name:             "indentless array",
		folder:           "indentless",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetIndentlessBlockSequence(true)
		},
	}.Run(t)
}

func TestIndentedToIndentless(t *testing.T) {
	formatTestCase{
		name:             "indented to indentless array",
		folder:           "indented_to_indentless",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetIndentlessBlockSequence(true)
		},
	}.Run(t)
}

func TestBlockScalar(t *testing.T) {
	formatTestCase{
		name:   "block scalar decoding and encoding",
		folder: "block_scalar",
		configureDecoder: func(dec *yaml.Decoder) {
			dec.SetScanBlockScalarAsLiteral(true)
		},
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetAssumeBlockAsLiteral(true)
		},
	}.Run(t)
}

func TestDropMergeTag(t *testing.T) {
	formatTestCase{
		name:             "drop merge tag",
		folder:           "drop_merge_tag",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetDropMergeTag(true)
		},
	}.Run(t)
}

func TestPadLineComments(t *testing.T) {
	formatTestCase{
		name:             "pad line comments",
		folder:           "pad_line_comments",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetPadLineComments(2)
		},
	}.Run(t)
}

func TestAltArrayIndent(t *testing.T) {
	formatTestCase{
		name:             "alternate array indent",
		folder:           "alt_array_indent",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetIndent(4)
			enc.SetArrayIndent(1)
		},
	}.Run(t)
}

func TestAltArrayIndentRoot(t *testing.T) {
	formatTestCase{
		name:             "alternate array indent (root)",
		folder:           "alt_array_indent_root",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetIndent(4)
			enc.SetArrayIndent(2)
			enc.SetIndentRootArray(true)
		},
	}.Run(t)
}

func TestFrontMatter(t *testing.T) {
	formatTestCase{
		name:             "frontmatter",
		folder:           "frontmatter_comments",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetExplicitDocumentStart(true)
		},
	}.Run(t)
}

func TestImplicitDocumentStartComments(t *testing.T) {
	formatTestCase{
		name:             "comment implicit document start",
		folder:           "comment_implicit_document_start",
		configureDecoder: noopDecoder,
		configureEncoder: func(enc *yaml.Encoder) {
			enc.SetExplicitDocumentStart(false)
		},
	}.Run(t)
}
