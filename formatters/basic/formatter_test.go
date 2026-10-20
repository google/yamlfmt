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
	"fmt"
	"strings"
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
			// Trailing whitespace on a content line used to force the emitter
			// to fall back to a double-quoted scalar (google/yamlfmt#86). Block
			// style is now preserved (the space_break heuristic no longer
			// disqualifies it), so the trailing space and the blank line both
			// survive in literal form. This also covers google/yamlfmt#280:
			// the retain_line_breaks placeholder would previously leak into
			// the quoted output here.
			name: "literal string with trailing space stays literal",
			input: "script: |\n" +
				"    when {\n" +
				"        deep\n" +
				"    }; \n" + // note trailing space
				"\n" +
				"    permit();\n",
			expect: "script: |\n" +
				"  when {\n" +
				"      deep\n" +
				"  }; \n" + // trailing space preserved
				"\n" +
				"  permit();\n",
		},
		{
			// A whitespace-only line inside a literal block scalar whose
			// width exceeds the block indent is content, not a blank line
			// (google/yamlfmt#86). The placeholder mechanism must round-trip
			// that content rather than collapsing it to an empty line.
			name: "literal string with whitespace-only content line",
			input: "script: |\n" +
				"  echo a\n" +
				"    \n" + // 4sp on a 2-indent block: content is "  "
				"  echo b\n",
			expect: "script: |\n" +
				"  echo a\n" +
				"    \n" +
				"  echo b\n",
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
		{
			// A blank line inside a > folded scalar is a paragraph break and
			// part of the scalar's value (YAML §8.1.3). The placeholder must
			// not be inserted there.
			name: "folded scalar blank line is content",
			input: "text: >\n" +
				"  para1\n" +
				"\n" +
				"  para2\n",
			expect: "text: >\n" +
				"  para1\n" +
				"\n" +
				"  para2\n",
		},
		{
			// A source-level blank line inside a multi-line "..." scalar
			// encodes a single \n in the value (YAML §7.3.1). The encoder
			// re-emits the value as a single-line double-quoted scalar with an
			// explicit \n; the round-tripped value must be unchanged.
			name: "double-quoted scalar blank line is content",
			input: "msg: \"line1\n" +
				"\n" +
				"  line2\"\n",
			expect: "msg: \"line1\\nline2\"\n",
		},
		{
			// Same for single-quoted; the encoder happens to keep the
			// multi-line layout, including the content blank.
			name: "single-quoted scalar blank line is content",
			input: "msg: 'line1\n" +
				"\n" +
				"  line2'\n",
			expect: "msg: 'line1\n" +
				"\n" +
				"  line2'\n",
		},
		{
			// retain_line_breaks_single chomps consecutive structural blanks
			// to one, but blank lines inside a literal scalar are content and
			// must not be chomped.
			name:   "chomp leaves scalar-content blanks alone",
			single: true,
			input: "script: |\n" +
				"  a\n" +
				"\n" +
				"\n" +
				"  b\n" +
				"x: 1\n" +
				"\n" +
				"\n" +
				"y: 2\n",
			expect: "script: |\n" +
				"  a\n" +
				"\n" +
				"\n" +
				"  b\n" +
				"x: 1\n" +
				"\n" +
				"y: 2\n",
		},
		{
			// A blank line between a sequence item and a less-indented
			// comment block. The placeholder must be inserted at the indent
			// of the *following* line: if it took the previous line's
			// (deeper) indent, the scanner would split it from the comment
			// block, attaching the comment to the next item instead and
			// emitting it at a different column. The mismatched column
			// reattaches differently on the next pass, so formatting would
			// only converge after a second run.
			name: "outdented comment after blank between sequence items",
			input: "items:\n" +
				"  - path: a\n" +
				"\n" +
				"# comment\n" +
				"  - path: b\n",
			expect: "items:\n" +
				"  - path: a\n" +
				"\n" +
				"    # comment\n" +
				"  - path: b\n",
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
			again, err := f.Format(got)
			require.NoError(t, err)
			require.Equal(t, string(got), string(again), "formatting is not idempotent")
		})
	}
}

// TestRetainLineBreaksLongLine is a regression test for a line longer than
// bufio.Scanner's default 64KiB token limit. The retain-line-breaks feature
// scans the document line by line; before the fix this failed with
// "bufio.Scanner: token too long" on real YAML carrying a large single-line
// value (e.g. SwiftPM .build manifests). See internal/linescan.
func TestRetainLineBreaksLongLine(t *testing.T) {
	longValue := strings.Repeat("x", 200*1024) // 200KiB, well past the 64KiB limit
	input := "a: 1\n\nbig: " + longValue + "\nb: 2\n"
	for _, single := range []bool{false, true} {
		t.Run(fmt.Sprintf("single=%t", single), func(t *testing.T) {
			f, err := factory.NewFormatter(map[string]any{
				"retain_line_breaks":        true,
				"retain_line_breaks_single": single,
			})
			require.NoError(t, err)
			got, err := f.Format([]byte(input))
			require.NoError(t, err)
			require.Contains(t, string(got), longValue, "the long single-line value must survive formatting")
		})
	}
}

func TestRetainLineBreaksFlowRestyledScalar(t *testing.T) {
	// With force_array_style=flow, a | block scalar inside an array is
	// re-emitted as a double-quoted flow scalar. The retain_line_breaks hotfix
	// skips placeholder insertion for blank lines inside the input scalar (the
	// scanner marks them as content), so no placeholder reaches the encoder and
	// none can leak into the restyled output.
	f, err := factory.NewFormatter(map[string]any{
		"retain_line_breaks": true,
		"force_array_style":  "flow",
	})
	require.NoError(t, err)
	in := "items:\n" +
		"  - |\n" +
		"    line1\n" +
		"\n" +
		"    line2\n"
	got, err := f.Format([]byte(in))
	require.NoError(t, err)
	require.NotContains(t, string(got), "magic___", "placeholder leaked into output")
	require.Equal(t, `items: ["line1\n\nline2\n"]`+"\n", string(got))
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

func TestRetainLineBreaksLiteralSentinelWithoutBlank(t *testing.T) {
	f, err := factory.NewFormatter(map[string]any{"retain_line_breaks": true})
	require.NoError(t, err)
	in := "message: |\n  #magic___^_^___line\n  hello\n"
	got, err := f.Format([]byte(in))
	require.NoError(t, err)
	require.Equal(t, in, string(got))
	again, err := f.Format(got)
	require.NoError(t, err)
	require.Equal(t, string(got), string(again))
}
