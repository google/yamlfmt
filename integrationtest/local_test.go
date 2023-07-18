package integrationtest

import (
	"flag"
	"testing"
)

// This suite contains tests that can easily be run on a local machine
// and operate entirely within temp directories that are created and
// destroyed per-test.

var (
	updateFlag = flag.Bool("update", false, "Whether to update the goldens.")
)

func TestPathArg(t *testing.T) {
	tempDirTestCase{
		Dir: "path_arg",
		// TODO: Path to command is the last thing to figure out before merge
		Command: "/home/braydon/go/bin/yamlfmt x.yaml",
		Update:  *updateFlag,
	}.Run(t)
}

func TestIncludeDocumentStart(t *testing.T) {
	tempDirTestCase{
		Dir:     "include_document_start",
		Command: "/home/braydon/go/bin/yamlfmt -formatter include_document_start=true x.yaml",
		Update:  *updateFlag,
	}.Run(t)
}
