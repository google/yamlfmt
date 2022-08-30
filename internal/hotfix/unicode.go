package hotfix

import (
	"errors"
	"regexp"

	"github.com/RageCage64/go-utf8-codepoint-converter/codepoint"
)

func ParseUnicodePoints(content []byte) []byte {
	if len(content) == 0 {
		return []byte{}
	}

	p := unicodeParser{
		buf: content,
		out: []byte{},
	}

	var err error
	for err != errEndOfBuffer {
		if p.peek() == '\\' {
			err = p.parseUTF8CodePoint()
			continue
		}

		p.write()
		err = p.next()
	}

	return p.out
}

type unicodeParser struct {
	buf []byte
	out []byte
	pos int
}

var (
	errInvalidCodePoint = errors.New("invalid UTF-8 codepoint sequence")
	errEndOfBuffer      = errors.New("end of buffer")
)

func (p *unicodeParser) peek() byte {
	return p.buf[p.pos]
}

func (p *unicodeParser) write() {
	p.out = append(p.out, p.peek())
}

func (p *unicodeParser) writeArbitrary(b []byte) {
	p.out = append(p.out, b...)
}

func (p *unicodeParser) next() error {
	p.pos++
	if p.pos == len(p.buf) {
		return errEndOfBuffer
	}
	return nil
}

func (p *unicodeParser) parseUTF8CodePoint() error {
	codepointBytes := []byte{}

	// Parse literal escape tokens while checking if this is a valid UTF-16 sequence
	if p.peek() != '\\' {
		return errInvalidCodePoint
	}
	codepointBytes = append(codepointBytes, p.peek())
	err := p.next()
	if err != nil {
		return err
	}
	if p.peek() != 'U' {
		return errInvalidCodePoint
	}
	codepointBytes = append(codepointBytes, p.peek())

	// We've detected a UTF-8 codepoint sequence. The library writes the UTF-16 sequence
	// i.e. \U0001F60A as 10 individual bytes. Our goal is to combine the 8
	// hexadecimal numbers we should see subsequently into the 4 byte values they
	// represent.
	isHex, err := regexp.Compile("[0-9A-F]")
	if err != nil {
		return err
	}

	for i := 0; i < 8; i++ {
		// Get a byte and confirm it is a hex digit.
		err = p.next()
		if err != nil {
			return err
		}
		hexDigit := p.peek()
		if !isHex.Match([]byte{hexDigit}) {
			return errInvalidCodePoint
		}
		codepointBytes = append(codepointBytes, hexDigit)
	}

	// Now that we have the codepoint, we'll represent it as a string
	// and pass it to the codepoint conversion library.
	utf8Bytes, err := codepoint.Convert(string(codepointBytes))
	if err != nil {
		return err
	}
	p.writeArbitrary(utf8Bytes)

	// Continue to the next byte for convenience to the caller.
	err = p.next()
	if err != nil {
		return err
	}

	return nil
}
