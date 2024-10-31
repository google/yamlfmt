package command

import (
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/internal/assert"
)

// This test asserts the proper behaviour for `line_ending` settings specified
// in formatter settings overriding the global configuration.
func TestLineEndingFormatterVsGloabl(t *testing.T) {
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
