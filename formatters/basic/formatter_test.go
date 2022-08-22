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

	"github.com/google/yamlfmt/formatters/basic"
)

func TestFormatterRetainsComments(t *testing.T) {
	f := &basic.BasicFormatter{Config: basic.DefaultConfig()}

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
