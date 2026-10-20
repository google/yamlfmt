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
	"strings"
	"testing"

	"github.com/google/yamlfmt/pkg/yaml"
	"github.com/stretchr/testify/require"
)

func TestRetainBlockScalarSeparators(t *testing.T) {
	for _, style := range []string{"|", "|-", "|+", ">", ">-", ">+", "|2-", ">+2"} {
		for _, single := range []bool{false, true} {
			for _, scan := range []bool{false, true} {
				for _, newline := range []string{"\n", "\r\n"} {
					t.Run(fmt.Sprintf("%s/single_%t/scan_%t/crlf_%t", style, single, scan, newline == "\r\n"), func(t *testing.T) {
						input := "block: " + style + "\n  first\n\n  second\n\nnext: value\n"
						input = strings.ReplaceAll(input, "\n", newline)
						f, err := factory.NewFormatter(map[string]any{"retain_line_breaks": true, "retain_line_breaks_single": single, "scan_folded_as_literal": scan, "line_ending": map[bool]string{false: "lf", true: "crlf"}[newline == "\r\n"]})
						require.NoError(t, err)
						got, err := f.Format([]byte(input))
						require.NoError(t, err)
						expected := strings.Replace(input, "block: "+style, "block: "+strings.ReplaceAll(style, "2", ""), 1)
						require.Equal(t, expected, string(got), "preserve both interior content and the separator after the scalar")
						again, err := f.Format(got)
						require.NoError(t, err)
						require.Equal(t, string(got), string(again))
						var before, after any
						require.NoError(t, yaml.Unmarshal([]byte(input), &before))
						require.NoError(t, yaml.Unmarshal(got, &after))
						require.Equal(t, before, after)
					})
				}
			}
		}
	}
}

func TestRetainLineBreaksBeyondTenMiB(t *testing.T) {
	input := "large: " + strings.Repeat("x", 10*1024*1024+1) + "\n\nnext: value\n"
	for _, single := range []bool{false, true} {
		t.Run(fmt.Sprintf("single_%t", single), func(t *testing.T) {
			f, err := factory.NewFormatter(map[string]any{"retain_line_breaks": true, "retain_line_breaks_single": single})
			require.NoError(t, err)
			got, err := f.Format([]byte(input))
			require.NoError(t, err)
			require.Equal(t, input, string(got))
			again, err := f.Format(got)
			require.NoError(t, err)
			require.Equal(t, string(got), string(again))
		})
	}
}
