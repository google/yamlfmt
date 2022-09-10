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
)

func newFormatter(config *basic.Config) *basic.BasicFormatter {
	formatter := &basic.BasicFormatter{Config: config}
	formatter.ConfigureFeaturesFromConfig()
	return formatter
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
	config.EmojiSupport = true
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
		desc   string
		input  string
		expect string
	}{
		{
			desc: "basic",
			input: `a:  1

b: 2`,
			expect: `a: 1

b: 2
`,
		},
		{
			desc: "multi-doc",
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
			desc: "literal string",
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
	}
	config := basic.DefaultConfig()
	config.RetainLineBreaks = true
	f := newFormatter(config)
	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			got, err := f.Format([]byte(c.input))
			if err != nil {
				t.Fatalf("expected formatting to pass, returned error: %v", err)
			}
			if string(got) != c.expect {
				t.Fatalf("didn't retain line breaks\nresult: %v\nexpect %s", string(got), c.expect)
			}
		})
	}
}
