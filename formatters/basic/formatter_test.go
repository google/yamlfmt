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

package basic_test

import (
	"testing"

	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/formatters/basic/features"
	"github.com/stretchr/testify/require"
)

var factory = basic.BasicFormatterFactory{}

func TestFormatter(t *testing.T) {
	testCases := []struct {
		name                     string
		config                   map[string]any
		badConfigErr             error
		formatErr                bool
		input                    string
		expect                   string
		skipLineEndNormalization bool
	}{
		{
			name:  "retains comments",
			input: `x: "y" # foo comment`,
		},
		{
			name: "parses multiple documents",
			input: `a:
---
b:`,
		},
		{
			name: "include document start",
			config: map[string]any{
				"include_document_start": true,
			},
			input: `a:`,
			expect: `---
a:`,
		},
		{
			name: "crlf line ending",
			config: map[string]any{
				"line_ending": "crlf",
			},
			input:  "a:\nb:\nc:\n",
			expect: "a:\r\nb:\r\nc:\r\n",

			skipLineEndNormalization: true,
		},
		{
			name: "lf line ending",
			config: map[string]any{
				"line_ending": "lf",
			},
			input:  "a:\r\nb:\r\nc:\r\n",
			expect: "a:\nb:\nc:\n",
		},
		{
			name:  "emoji support",
			input: `a: 😊`,
		},
		{
			name: "scan folded as literal",
			config: map[string]any{
				"scan_folded_as_literal": true,
				"retain_line_breaks":     true,
			},
			input: `a: >
  multiline
  folded
  scalar`,
		},
		{
			name: "indentless arrays",
			config: map[string]any{
				"indentless_arrays": true,
			},
			input: `a:
  - 1
  - 2`,
			expect: `a:
- 1
- 2`,
		},
		{
			name: "drop merge tag",
			config: map[string]any{
				"drop_merge_tag": true,
			},
			input: `a: &a
  b:
    <<: *a`,
		},
		{
			name: "pad line comments",
			config: map[string]any{
				"pad_line_comments": 2,
			},
			input:  `a: 1 # line comment`,
			expect: `a: 1  # line comment`,
		},
		{
			name: "trim trailing whitespace",
			config: map[string]any{
				"trim_trailing_whitespace": true,
			},
			input: `a: 1
b: 2    `,
			expect: `a: 1
b: 2`,
		},
		{
			name: "eof newline",
			config: map[string]any{
				"eof_newline": true,
			},
			input: `a: 1
b: 2`,
			expect: `a: 1
b: 2
`,
		},
		{
			name: "strip directives",
			config: map[string]any{
				"strip_directives": true,
			},
			input: "%YAML:1.0\na: 1",
		},
		{
			name: "array indent",
			config: map[string]any{
				"array_indent": 2,
			},
			input: `a:
    - 1
    - 2`,
			expect: `a:
  - 1
  - 2`,
		},
		{
			name: "indent root array",
			config: map[string]any{
				"indent_root_array": true,
				"array_indent":      2,
			},
			input: `  - 1
  - 2`,
			expect: `  - 1
  - 2`,
		},
		{
			name: "force block sequence style",
			config: map[string]any{
				"force_array_style": "block",
			},
			input: `a:
  - 1
  - 2
b: [1, 2]`,
			expect: `a:
  - 1
  - 2
b:
  - 1
  - 2`,
		},
		{
			name: "force flow sequence style",
			config: map[string]any{
				"force_array_style": "flow",
			},
			input: `a:
  - 1
  - 2
b: [1, 2]`,
			expect: `a: [1, 2]
b: [1, 2]`,
		},
		{
			name: "invalid sequence style",
			config: map[string]any{
				"force_array_style": "invalid",
			},
			badConfigErr: features.ErrUnrecognizedSequenceStyle,
		},
		{
			name: "force single quote style",
			config: map[string]any{
				"force_quote_style": "single",
			},
			input:  `a: ['hi', "hello"]`,
			expect: `a: ['hi', 'hello']`,
		},
		{
			name: "force double quote style",
			config: map[string]any{
				"force_quote_style": "double",
			},
			input:  `a: ['hi', "hello"]`,
			expect: `a: ["hi", "hello"]`,
		},
		{
			name: "invalid quote style",
			config: map[string]any{
				"force_quote_style": "blah",
			},
			badConfigErr: features.ErrUnrecognizedQuoteStyle,
		},
		{
			name: "alias key correction",
			input: `alias: &a something
map:
  *a: 1`,
			expect: `alias: &a something
map:
  *a : 1`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := factory.NewFormatter(tc.config)
			if tc.badConfigErr != nil {
				require.ErrorIs(t, err, tc.badConfigErr)
				return
			}
			require.NoError(t, err)

			result, err := f.Format([]byte(tc.input))
			if tc.formatErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected := tc.expect
			actual := string(result)
			if tc.expect == "" {
				expected = tc.input
			}

			if !tc.skipLineEndNormalization {
				// Having to always include the newline in the expected
				// data in the test table is cringe unless that's something
				// I actually want to test
				expected = stripTrailingNewline(expected)
				actual = stripTrailingNewline(actual)
			}
			require.Equal(t, expected, actual)
		})
	}
}

func TestRetainLineBreaks(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect string
		single bool
	}{
		{
			name: "basic",
			input: `a:  1

b: 2`,
			expect: `a: 1

b: 2
`,
		},
		{
			name: "multi-doc",
			input: `a:  1

# tail comment
---
b: 2`,
			expect: `a: 1

# tail comment
---
b: 2
`,
		},
		{
			name: "literal string",
			input: `a:  1

shell: |
  #!/usr/bin/env bash

  # hello, world
    # bye
  echo "hello, world"
`,
			expect: `a: 1

shell: |
  #!/usr/bin/env bash

  # hello, world
    # bye
  echo "hello, world"
`,
		},
		{
			name: "multi level nested literal string",
			input: `a:  1
x:
  y:
    shell: |
      #!/usr/bin/env bash

        # bye
      echo "hello, world"`,
			expect: `a: 1
x:
  y:
    shell: |
      #!/usr/bin/env bash

        # bye
      echo "hello, world"
`,
		},
		{
			name:   "retain single line break",
			single: true,
			input: `a: 1

b: 2

c: 3
`,
			expect: `a: 1

b: 2

c: 3
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := factory.NewFormatter(map[string]any{
				"retain_line_breaks":        true,
				"retain_line_breaks_single": tc.single,
			})
			require.NoError(t, err)
			got, err := f.Format([]byte(tc.input))
			require.NoError(t, err)
			require.Equal(t, tc.expect, string(got))
		})
	}
}

func stripTrailingNewline(s string) string {
	// strip trailing \n or \r\n characters
	if len(s) > 0 {
		if s[len(s)-1] == '\n' {
			s = s[:len(s)-1]
		} else if len(s) > 1 && s[len(s)-2] == '\r' {
			s = s[:len(s)-2]
		}
	}
	return s
}
