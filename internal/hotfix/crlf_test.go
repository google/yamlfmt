package hotfix_test

import (
	"testing"

	"github.com/google/yamlfmt/internal/hotfix"
)

func TestStripCRBytes(t *testing.T) {
	crlfContent := []byte("two\r\nlines\r\n")
	lfContent := hotfix.StripCRBytes(crlfContent)
	count := countByte(lfContent, '\r')
	if count != 0 {
		t.Fatalf("Found %d CR (decimal 13) bytes in %v", count, lfContent)
	}
}

func TestWriteCRLF(t *testing.T) {
	lfContent := []byte("two\nlines\n")
	crlfContent := hotfix.WriteCRLFBytes(lfContent)
	countCR := countByte(crlfContent, '\r')
	countLF := countByte(crlfContent, '\n')
	if countCR != countLF {
		t.Fatalf("Found %d CR (decimal 13) and %d LF (decimal 10) bytes in %v", countCR, countLF, crlfContent)
	}
}

func countByte(content []byte, searchByte byte) int {
	count := 0
	for _, b := range content {
		if b == searchByte {
			count++
		}
	}
	return count
}
