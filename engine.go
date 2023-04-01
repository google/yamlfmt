package yamlfmt

import (
	"fmt"
	"os"

	"github.com/RageCage64/multilinediff"
	"github.com/google/yamlfmt/internal/collections"
)

type Engine interface {
	FormatContent(content []byte) ([]byte, error)
	Format(paths []string) error
	Lint(paths []string) (*EngineOutput, error)
	DryRun(paths []string) (*EngineOutput, error)
}

type EngineOutput struct {
	Message string
	Files   FileDiffs
	Quiet   bool
}

func (eo *EngineOutput) String() string {
	var msg string
	if eo.Message != "" {
		msg += fmt.Sprintf("%s\n\n", eo.Message)
	}
	if eo.Quiet {
		msg += eo.Files.StrOutputQuiet()
	} else {
		msg += fmt.Sprintf("%s\n", eo.Files.StrOutput())
	}
	return msg
}

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
	return fd.Path + "\n"
}

func (fd *FileDiff) Apply() error {
	return os.WriteFile(fd.Path, []byte(fd.Diff.Formatted), 0644)
}

type FileDiffs map[string]*FileDiff

func (fds FileDiffs) Add(diff *FileDiff) error {
	if _, ok := fds[diff.Path]; ok {
		return fmt.Errorf("a diff for %s already exists", diff.Path)
	}

	fds[diff.Path] = diff
	return nil
}

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
	applyErrs := make(collections.Errors, len(fds))
	i := 0
	for _, diff := range fds {
		applyErrs[i] = diff.Apply()
		i++
	}
	return applyErrs.Combine()
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
