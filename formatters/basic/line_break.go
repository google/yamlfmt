package basic

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const lineBreakPlaceholder = "#magic___^_^___line"

// retainLineBreaks keeps the line breaks.
//
// The basic idea is to insert/remove placeholder comments in the yaml document before and after the format process.
func retainLineBreaks(in io.Reader, formatter func(io.Reader) (io.Reader, error)) ([]byte, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(in)
	var padding paddinger
	for scanner.Scan() {
		txt := scanner.Text()
		padding.adjust(indentLen(txt))
		if strings.TrimSpace(txt) == "" { // line break or empty space line.
			buf.WriteString(padding.String()) // prepend some padding incase literal multiline strings.
			buf.WriteString(lineBreakPlaceholder)
			buf.WriteByte('\n')
			continue
		}
		buf.WriteString(txt)
		buf.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	formatted, err := formatter(&buf)
	if err != nil {
		return nil, err
	}

	buf.Reset() // reuse

	scanner = bufio.NewScanner(formatted)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.TrimSpace(txt) == "" {
			// The basic yaml lib inserts newline when there is a comment(either placeholder or by user)
			// followed by optional line breaks and a `---` multi-documents.
			// To fix it, the empty line could only be inserted by us.
			continue
		}
		if strings.HasPrefix(strings.TrimLeft(txt, " "), lineBreakPlaceholder) {
			buf.WriteByte('\n')
			continue
		}
		buf.WriteString(txt)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), scanner.Err()
}

func indentLen(txt string) int {
	var cnt int
	for i := 0; i < len(txt); i++ {
		if txt[i] != ' ' { // Yaml only allows space to indent.
			break
		}
		cnt++
	}
	return cnt
}

type paddinger struct {
	strings.Builder
}

func (p *paddinger) adjust(size int) {
	// Grows if the given size is larger than us and always return the max padding.
	for diff := size - p.Len(); diff > 0; diff-- {
		p.WriteByte(' ')
	}
}
