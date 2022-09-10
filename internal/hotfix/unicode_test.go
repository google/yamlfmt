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

package hotfix_test

import (
	"testing"

	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/internal/hotfix"
)

func TestParseEmoji(t *testing.T) {
	testCases := []struct {
		name        string
		yamlStr     string
		expectedStr string
	}{
		{
			name:        "parses emoji",
			yamlStr:     "a: 😂\n",
			expectedStr: "a: \"😂\"\n",
		},
		{
			name:        "parses multiple emoji",
			yamlStr:     "a: 😼 👑\n",
			expectedStr: "a: \"😼 👑\"\n",
		},
	}

	f := &basic.BasicFormatter{Config: basic.DefaultConfig()}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formattedBefore, err := f.Format([]byte(tc.yamlStr))
			if err != nil {
				t.Fatalf("yaml failed to parse: %v", err)
			}
			formattedAfter, err := hotfix.ParseUnicodePoints(formattedBefore)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			formattedStr := string(formattedAfter)
			if formattedStr != tc.expectedStr {
				t.Fatalf("parsed string does not match: \nexpected: %s\ngot: %s", tc.expectedStr, string(formattedStr))
			}
		})
	}
}
