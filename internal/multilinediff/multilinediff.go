package multilinediff

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Get the diff between two strings.
func Diff(a, b, lineSep string) (string, int) {
	reporter := Reporter{LineSep: lineSep}
	cmp.Diff(
		a, b,
		cmpopts.AcyclicTransformer("multiline", func(s string) []string {
			return strings.Split(s, lineSep)
		}),
		cmp.Reporter(&reporter),
	)
	return reporter.String(), reporter.DiffCount
}

type diffType int

const (
	diffTypeEqual diffType = iota
	diffTypeChange
	diffTypeAdd
)

type diffLine struct {
	diff diffType
	old  string
	new  string
}

func (l diffLine) toLine(length int) string {
	line := ""

	switch l.diff {
	case diffTypeChange:
		line += "- "
	case diffTypeAdd:
		line += "+ "
	default:
		line += "  "
	}

	line += l.old

	for i := 0; i < length-len(l.old); i++ {
		line += " "
	}

	line += "  "

	line += l.new

	return line
}

// A pretty reporter to pass into cmp.Diff using the cmd.Reporter function.
type Reporter struct {
	LineSep   string
	DiffCount int

	path  cmp.Path
	lines []diffLine
}

func (r *Reporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *Reporter) Report(rs cmp.Result) {
	line := diffLine{}
	vOld, vNew := r.path.Last().Values()
	if !rs.Equal() {
		r.DiffCount++
		if vOld.IsValid() {
			line.diff = diffTypeChange
			line.old = fmt.Sprintf("%+v", vOld)
		}
		if vNew.IsValid() {
			if line.diff == diffTypeEqual {
				line.diff = diffTypeAdd
			}
			line.new = fmt.Sprintf("%+v", vNew)
		}
	} else {
		line.old = fmt.Sprintf("%+v", vOld)
		line.new = fmt.Sprintf("%+v", vOld)
	}
	r.lines = append(r.lines, line)
}

func (r *Reporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *Reporter) String() string {
	maxLen := 0
	for _, l := range r.lines {
		if len(l.old) > maxLen {
			maxLen = len(l.old)
		}
	}

	diffLines := []string{}
	for _, l := range r.lines {
		diffLines = append(diffLines, l.toLine(maxLen))
	}

	return strings.Join(diffLines, r.LineSep)
}
