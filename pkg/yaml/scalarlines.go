package yaml

// ScalarLineRange marks a half-open range of 0-indexed source lines that
// belong to a single scalar token. StartLine is the line of the scalar's
// first character (for |/> block scalars this is the indicator line), and
// EndLine excludes trailing block-scalar breaks discarded by chomping.
type ScalarLineRange struct {
	StartLine int
	EndLine   int
}

// ScanMultilineScalarRanges runs the low-level scanner over input and
// returns the source line range of every scalar token that spans more than
// one line. It is used by the retain-line-break hotfix to avoid rewriting
// blank lines that are scalar content rather than structural whitespace.
//
// On scan error it returns whatever ranges were collected before the error:
// callers prefer best-effort protection over failing the whole format pass
// on a file the real decoder may still accept.
func ScanMultilineScalarRanges(input []byte) []ScalarLineRange {
	var parser yaml_parser_t
	if !yaml_parser_initialize(&parser) {
		return nil
	}
	yaml_parser_set_input_string(&parser, input)

	var ranges []ScalarLineRange
	var tok yaml_token_t
	for {
		if !yaml_parser_scan(&parser, &tok) {
			return ranges
		}
		if tok.typ == yaml_STREAM_END_TOKEN {
			return ranges
		}
		if tok.typ != yaml_SCALAR_TOKEN {
			continue
		}
		if tok.end_mark.line <= tok.start_mark.line {
			continue
		}
		end := tok.end_mark.line
		if tok.style == yaml_LITERAL_SCALAR_STYLE || tok.style == yaml_FOLDED_SCALAR_STYLE {
			end = tok.scalar_content_end_mark.line
		}
		ranges = append(ranges, ScalarLineRange{
			StartLine: tok.start_mark.line,
			EndLine:   end,
		})
	}
}
