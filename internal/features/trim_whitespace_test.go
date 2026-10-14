// Copyright 2026 Google LLC
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

package features

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/yamlfmt/pkg/yaml"
	"github.com/stretchr/testify/require"
)

func TestTrimTrailingWhitespacePreservesChomping(t *testing.T) {
	for _, style := range []string{"|", "|+", "|-", ">", ">+", ">-"} {
		for _, inputLineEnding := range []string{"\n", "\r\n", "\r"} {
			for _, outputLineEnding := range []string{"\n", "\r\n"} {
				for _, trailingBreaks := range []int{0, 1, 2} {
					name := fmt.Sprintf("%s/input=%q/output=%q/breaks=%d", style, inputLineEnding, outputLineEnding, trailingBreaks)
					t.Run(name, func(t *testing.T) {
						input := "scalar: " + style + inputLineEnding + "  foo" + strings.Repeat(inputLineEnding, trailingBreaks)
						_, output, err := MakeFeatureTrimTrailingWhitespace(outputLineEnding).BeforeAction(context.Background(), []byte(input))
						require.NoError(t, err)
						var before, after map[string]any
						require.NoError(t, yaml.Unmarshal([]byte(input), &before))
						require.NoError(t, yaml.Unmarshal(output, &after))
						require.Equal(t, before, after)
					})
				}
			}
		}
	}
}
