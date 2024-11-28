// Copyright 2024 Google LLC
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

package command

import (
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/internal/assert"
)

// This test asserts the proper behaviour for `line_ending` settings specified
// in formatter settings overriding the global configuration.
func TestLineEndingFormatterVsGlobal(t *testing.T) {
	c := &Command{
		Config: &Config{
			LineEnding: "lf",
			FormatterConfig: &FormatterConfig{
				FormatterSettings: map[string]any{
					"line_ending": yamlfmt.LineBreakStyleLF,
				},
			},
		},
		Registry: yamlfmt.NewFormatterRegistry(&basic.BasicFormatterFactory{}),
	}

	f, err := c.getFormatter()
	assert.NilErr(t, err)
	configMap, err := f.ConfigMap()
	assert.NilErr(t, err)
	formatterLineEnding := configMap["line_ending"].(yamlfmt.LineBreakStyle)
	assert.Assert(t, formatterLineEnding == yamlfmt.LineBreakStyleLF, "expected formatter's line ending to be lf")
}
