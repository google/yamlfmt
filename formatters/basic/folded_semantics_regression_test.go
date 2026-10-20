// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package basic_test

import (
	"fmt"
	"testing"

	"github.com/google/yamlfmt/pkg/yaml"
	"github.com/stretchr/testify/require"
)

func TestFoldedScalarValueAndIdempotence(t *testing.T) {
	inputs := map[string]string{
		"whitespace_line": "message: >-\n  first line\n    \n  next line\n\nnext: value\n",
		"keep_chomping":   "message: >+\n  first paragraph\n\n  second paragraph\n\nnext: value\n",
		"trailing_space":  "message: >+\n  first line   \n  next line\n\nnext: value\n",
	}
	for name, input := range inputs {
		for _, retain := range []string{"retain_line_breaks", "retain_line_breaks_single"} {
			for _, scan := range []bool{false, true} {
				for _, width := range []int{0, 40, 80} {
					t.Run(fmt.Sprintf("%s/%s/scan_%t/width_%d", name, retain, scan, width), func(t *testing.T) {
						f, err := factory.NewFormatter(map[string]any{retain: true, "scan_folded_as_literal": scan, "max_line_length": width, "eof_newline": true})
						require.NoError(t, err)
						got, err := f.Format([]byte(input))
						require.NoError(t, err)
						var before, after any
						require.NoError(t, yaml.Unmarshal([]byte(input), &before))
						require.NoError(t, yaml.Unmarshal(got, &after))
						require.Equal(t, before, after, "formatting changed the YAML value")
						again, err := f.Format(got)
						require.NoError(t, err)
						require.Equal(t, string(got), string(again), "formatting is not idempotent")
					})
				}
			}
		}
	}
}
