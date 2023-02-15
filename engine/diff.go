package engine

import (
	"errors"
	"fmt"
	"os"

	"github.com/RageCage64/multilinediff"
)

type FormatDiff struct {
	Original  string
	Formatted string
	LineSep   string
}

func (d *FormatDiff) MultilineDiff() (string, int) {
	return multilinediff.Diff(d.Original, d.Formatted, d.LineSep)
}

func (d *FormatDiff) Changed() bool {
	return d.Original != d.Formatted
}

type FileDiff struct {
	Path string
	Diff *FormatDiff
}

func (fd *FileDiff) StrOutput() string {
	diffStr, _ := fd.Diff.MultilineDiff()
	return fmt.Sprintf("%s:\n%s\n", fd.Path, diffStr)
}

func (fd *FileDiff) StrOutputQuiet() string {
	return fd.Path
}

func (fd *FileDiff) Apply() error {
	return os.WriteFile(fd.Path, []byte(fd.Diff.Formatted), 0644)
}

type FileDiffs []*FileDiff

func (fds FileDiffs) StrOutput() string {
	result := ""
	for _, fd := range fds {
		if fd.Diff.Changed() {
			result += fd.StrOutput()
		}
	}
	return result
}

func (fds FileDiffs) StrOutputQuiet() string {
	result := ""
	for _, fd := range fds {
		if fd.Diff.Changed() {
			result += fd.StrOutputQuiet()
		}
	}
	return result
}

func (fds FileDiffs) ApplyAll() error {
	applyErrs := make([]error, len(fds))
	for i, diff := range fds {
		applyErrs[i] = diff.Apply()
	}
	return errors.Join(applyErrs...)
}

func (fds FileDiffs) ChangedCount() int {
	changed := 0
	for _, fd := range fds {
		if fd.Diff.Changed() {
			changed++
		}
	}
	return changed
}
