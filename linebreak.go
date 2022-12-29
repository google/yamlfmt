package yamlfmt

import "fmt"

type LineBreakStyle string

const (
	LineBreakStyleLF   LineBreakStyle = "lf"
	LineBreakStyleCRLF LineBreakStyle = "crlf"
)

type UnsupportedLineBreakError struct {
	style LineBreakStyle
}

func (e UnsupportedLineBreakError) Error() string {
	return fmt.Sprintf("unsupported line break style %s, see package documentation for supported styles", e.style)
}

func (s LineBreakStyle) Separator() (string, error) {
	switch s {
	case LineBreakStyleLF:
		return "\n", nil
	case LineBreakStyleCRLF:
		return "\r\n", nil
	}
	return "", UnsupportedLineBreakError{style: s}
}
