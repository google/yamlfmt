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
	"strings"
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/internal/assert"
)

func newFormatter(config *basic.Config) *basic.BasicFormatter {
	return &basic.BasicFormatter{
		Config:   config,
		Features: basic.ConfigureFeaturesFromConfig(config),
	}
}

func TestFormatterRetainsComments(t *testing.T) {
	f := newFormatter(basic.DefaultConfig())

	yaml := `x: "y" # foo comment`

	s, err := f.Format([]byte(yaml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	if !strings.Contains(string(s), "#") {
		t.Fatal("comment was stripped away")
	}
}

func TestFormatterPreservesKeyOrder(t *testing.T) {
	f := &basic.BasicFormatter{Config: basic.DefaultConfig()}

	yaml := `
b:
a:`

	s, err := f.Format([]byte(yaml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	unmarshalledStr := string(s)
	bPos := strings.Index(unmarshalledStr, "b")
	aPos := strings.Index(unmarshalledStr, "a")
	if bPos > aPos {
		t.Fatalf("keys were reordered:\n%s", s)
	}
}

func TestFormatterParsesMultipleDocuments(t *testing.T) {
	f := &basic.BasicFormatter{Config: basic.DefaultConfig()}

	yaml := `b:
---
a:
`
	s, err := f.Format([]byte(yaml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	if len(s) != len([]byte(yaml)) {
		t.Fatalf("expected yaml not to change, result: %s", string(s))
	}
}

func TestWithDocumentStart(t *testing.T) {
	config := basic.DefaultConfig()
	config.IncludeDocumentStart = true
	f := newFormatter(config)

	yaml := "a:"
	s, err := f.Format([]byte(yaml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	if strings.Index(string(s), "---\n") != 0 {
		t.Fatalf("expected document start to be included, result was: %s", string(s))
	}
}

func TestCRLFLineEnding(t *testing.T) {
	config := basic.DefaultConfig()
	config.LineEnding = yamlfmt.LineBreakStyleCRLF
	f := newFormatter(config)

	yaml := "# comment\r\na:\r\n"
	result, err := f.Format([]byte(yaml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	if string(yaml) != string(result) {
		t.Fatalf("didn't write CRLF properly in result: %v", result)
	}
}

func TestEmojiSupport(t *testing.T) {
	config := basic.DefaultConfig()
	f := newFormatter(config)

	yaml := "a: 😊"
	result, err := f.Format([]byte(yaml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := string(result)
	if !strings.Contains(resultStr, "😊") {
		t.Fatalf("expected string to contain 😊, got: %s", resultStr)
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
			config := basic.DefaultConfig()
			config.RetainLineBreaks = true
			config.RetainLineBreaksSingle = tc.single
			f := newFormatter(config)
			got, err := f.Format([]byte(tc.input))
			if err != nil {
				t.Fatalf("expected formatting to pass, returned error: %v", err)
			}
			if string(got) != tc.expect {
				t.Fatalf("didn't retain line breaks\nresult: %v\nexpect %s", string(got), tc.expect)
			}
		})
	}
}

func TestScanFoldedAsLiteral(t *testing.T) {
	config := basic.DefaultConfig()
	config.ScanFoldedAsLiteral = true
	f := newFormatter(config)

	yml := `a: >
  multiline
  folded
  scalar`
	lines := len(strings.Split(yml, "\n"))
	result, err := f.Format([]byte(yml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := string(result)
	resultLines := len(strings.Split(resultStr, "\n"))
	if resultLines == lines {
		t.Fatalf("expected string to be %d lines, was %d", lines, resultLines)
	}
}

func TestIndentlessArrays(t *testing.T) {
	config := basic.DefaultConfig()
	config.IndentlessArrays = true
	f := newFormatter(config)

	yml := `a:
- 1
- 2
`
	result, err := f.Format([]byte(yml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := string(result)
	if resultStr != yml {
		t.Fatalf("expected:\n%s\ngot:\n%s", yml, resultStr)
	}
}

func TestDropMergeTag(t *testing.T) {
	config := basic.DefaultConfig()
	config.DropMergeTag = true
	f := newFormatter(config)

	yml := `a: &a
b:
  <<: *a`

	result, err := f.Format([]byte(yml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := string(result)
	if strings.Contains(resultStr, "!!merge") {
		t.Fatalf("expected formatted result to drop merge tag, was found:\n%s", resultStr)
	}
}

func TestPadLineComments(t *testing.T) {
	config := basic.DefaultConfig()
	config.PadLineComments = 2
	f := newFormatter(config)

	yml := "a: 1 # line comment"
	expectedStr := "a: 1  # line comment"

	result, err := f.Format([]byte(yml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := strings.TrimSuffix(string(result), "\n")
	if resultStr != expectedStr {
		t.Fatalf("expected: '%s', got: '%s'", expectedStr, resultStr)
	}
}

func TestTrimTrailingWhitespace(t *testing.T) {
	config := basic.DefaultConfig()
	config.TrimTrailingWhitespace = true
	f := newFormatter(config)

	yml := `a: 1
b: 2    `
	expectedYml := `a: 1
b: 2`

	result, err := f.Format([]byte(yml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := strings.TrimSuffix(string(result), "\n")
	if resultStr != expectedYml {
		t.Fatalf("expected: '%s', got: '%s'", expectedYml, resultStr)
	}
}

func TestEOFNewline(t *testing.T) {
	config := basic.DefaultConfig()
	config.RetainLineBreaks = false
	config.EOFNewline = true
	f := newFormatter(config)

	yml := `a: 1
b: 2`
	expectedYml := `a: 1
b: 2
`

	result, err := f.Format([]byte(yml))
	if err != nil {
		t.Fatalf("expected formatting to pass, returned error: %v", err)
	}
	resultStr := string(result)
	if resultStr != expectedYml {
		t.Fatalf("expected: '%s', got: '%s'", expectedYml, resultStr)
	}
}

func TestStripDirectives(t *testing.T) {
	config := basic.DefaultConfig()
	config.StripDirectives = true
	f := newFormatter(config)

	yml := "%YAML:1.0"

	_, err := f.Format([]byte(yml))
	assert.NilErr(t, err)
}

func TestArrayIndent(t *testing.T) {
	config := basic.DefaultConfig()
	config.ArrayIndent = 1
	f := newFormatter(config)

	yml := `a:
  - 1
  - 2
`
	expectedYml := `a:
 - 1
 - 2
`

	result, err := f.Format([]byte(yml))
	assert.NilErr(t, err)
	resultStr := string(result)
	if resultStr != expectedYml {
		t.Fatalf("expected: '%s', got: '%s'", expectedYml, result)
	}
}

func TestIndentRootArray(t *testing.T) {
	config := basic.DefaultConfig()
	config.IndentRootArray = true
	f := newFormatter(config)

	yml := "- 1\n"
	expectedYml := "  - 1\n"

	result, err := f.Format([]byte(yml))
	assert.NilErr(t, err)
	resultStr := string(result)
	if resultStr != expectedYml {
		t.Fatalf("expected: '%s', got: '%s'", expectedYml, result)
	}
}
