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

package diff_test

import (
	"fmt"
	"testing"

	"github.com/google/yamlfmt/internal/diff"
)

func TestMultilineDiff(t *testing.T) {
	testCases := []struct {
		name   string
		before string
		after  string
	}{
		{
			name: "multiple line diff",
			before: `this
is
a
first input`,
			after: `this
is
the
second input`,
		},
		{
			name: "add blank line",
			before: `foo
bar`,
			after: `foo

bar`,
		},
		{
			name: "remove blank line",
			before: `foo
			
bar`,
			after: `foo
bar`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diffStr := diff.MultilineStringDiff(tc.before, tc.after)
			if diffStr == "" {
				t.Log("there should have been a diff")
				t.Fail()
			}
			fmt.Println(diffStr)
		})
	}
}

func TestMultilineDiffNoDiff(t *testing.T) {
	before := "content"

	diffStr := diff.MultilineStringDiff(before, before)
	if diffStr != "" {
		t.Fatalf("diff output should be empty with no diff, output contained:\n%s", diffStr)
	}
}
