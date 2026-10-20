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

// The features in this file are to retain line breaks.
// The basic idea is to insert/remove placeholder comments in the yaml document before and after the format process.
//
// A blank line inside a multi-line scalar (|, >, "...", '...', or plain) is
// content, not structural whitespace. The original implementation inserted the
// placeholder for every blank line, which corrupted such scalars
// (google/yamlfmt#280, #86). To avoid that, both actions first run the YAML
// scanner over their input and skip the rewrite for any line that falls inside
// a multi-line scalar token. The scanner is authoritative, so this covers all
// scalar styles without lexical heuristics.

package hotfix

import (
	"bytes"
	"context"
	"strings"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/linescan"
	"github.com/google/yamlfmt/pkg/yaml"
)

const lineBreakPlaceholder = "#magic___^_^___line"

// scalarContentLines protects blank lines that contribute to scalar values.
// Scanner ranges exclude trailing block-scalar breaks discarded by chomping,
// so structural separators can still be retained using placeholders.
func scalarContentLines(src []byte) map[int]struct{} {
	in := make(map[int]struct{})
	for _, r := range yaml.ScanMultilineScalarRanges(src) {
		for l := r.StartLine + 1; l < r.EndLine; l++ {
			in[l] = struct{}{}
		}
	}
	return in
}

// indentOf returns the leading-space prefix of txt. YAML only allows spaces
// for indentation, so tabs are not considered.
func indentOf(txt string) string {
	for i := 0; i < len(txt); i++ {
		if txt[i] != ' ' {
			return txt[:i]
		}
	}
	return txt
}

func MakeFeatureRetainLineBreak(linebreakStr string, chomp bool) yamlfmt.Feature {
	return yamlfmt.Feature{
		Name:         "Retain Line Breaks",
		BeforeAction: replaceLineBreakFeature(linebreakStr, chomp),
		AfterAction:  restoreLineBreakFeature(linebreakStr),
	}
}

func replaceLineBreakFeature(newlineStr string, chomp bool) yamlfmt.FeatureFunc {
	return func(ctx context.Context, content []byte) (context.Context, []byte, error) {
		inScalar := scalarContentLines(content)
		var buf bytes.Buffer
		scanner := linescan.NewScanner(content)
		var lineNo int
		// Placeholders are written at the indent of the line that *follows*
		// the blank run, not the line that precedes it. The scanner uses
		// comment indentation to decide attachment, so a placeholder at the
		// previous line's (deeper) indent splits it from a less-indented
		// comment block that follows, attaching that block to the next node
		// instead. The mismatched column then re-attaches differently on the
		// next pass, so formatting only converges after a second run.
		// Matching the following line keeps placeholder + following content as
		// one unit and the placeholder transparent to attachment.
		var pendingPlaceholders int
		flushPlaceholders := func(indent string) {
			for ; pendingPlaceholders > 0; pendingPlaceholders-- {
				buf.WriteString(indent)
				buf.WriteString(lineBreakPlaceholder)
				buf.WriteString(newlineStr)
			}
		}
		for scanner.Scan() {
			txt := scanner.Text()
			_, protected := inScalar[lineNo]
			lineNo++
			if strings.TrimRight(txt, " \t") == "" {
				if protected {
					// Blank (or whitespace-only) line that is scalar content:
					// pass through verbatim so the decoder sees the original
					// value. The encoder preserves it natively. This is not a
					// structural blank, so don't count it as pending.
					buf.WriteString(txt)
					buf.WriteString(newlineStr)
					continue
				}
				if chomp && pendingPlaceholders > 0 {
					continue
				}
				pendingPlaceholders++
				continue
			}
			flushPlaceholders(indentOf(txt))
			buf.WriteString(txt)
			buf.WriteString(newlineStr)
		}
		// Trailing structural blanks have no following line; column zero is
		// transparent (anything outdented attaches at document scope, same as
		// the original blank).
		flushPlaceholders("")
		return ctx, buf.Bytes(), scanner.Err()
	}
}

func restoreLineBreakFeature(newlineStr string) yamlfmt.FeatureFunc {
	return func(ctx context.Context, content []byte) (context.Context, []byte, error) {
		// Always scan output: literal user content may match the sentinel even
		// when the input contained no blank line inside a scalar.
		inScalar := scalarContentLines(content)
		var buf bytes.Buffer
		scanner := linescan.NewScanner(content)
		var lineNo int
		for scanner.Scan() {
			txt := scanner.Text()
			_, protected := inScalar[lineNo]
			lineNo++
			if txt == "" {
				if protected {
					buf.WriteString(newlineStr)
					continue
				}
				// Structural blank in encoder output that we did not request
				// via a placeholder: drop it (matches the original hotfix's
				// handling of the comment-before-`---` quirk).
				continue
			}
			if strings.HasPrefix(strings.TrimLeft(txt, " "), lineBreakPlaceholder) {
				if protected {
					// The placeholder text appears inside a scalar in the
					// output. The BeforeAction never inserts placeholders into
					// scalar content, so this can only be literal user content
					// that happens to match the sentinel. Round-trip it.
					buf.WriteString(txt)
					buf.WriteString(newlineStr)
					continue
				}
				buf.WriteString(newlineStr)
				continue
			}
			buf.WriteString(txt)
			buf.WriteString(newlineStr)
		}
		return nil, buf.Bytes(), scanner.Err()
	}
}
