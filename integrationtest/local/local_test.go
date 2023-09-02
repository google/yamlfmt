//go:build integration_test

package local_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/yamlfmt/integrationtest/local"
)

var (
	updateFlag *bool = flag.Bool("update", false, "Whether to update the goldens.")
	yamlfmtBin string
)

func init() {
	yamlfmtBinVar := os.Getenv("YAMLFMT_BIN")
	if yamlfmtBinVar == "" {
		fmt.Println("Must provide a YAMLFMT_BIN environment variable.")
		os.Exit(1)
	}
	yamlfmtBin = yamlfmtBinVar
}

func TestPathArg(t *testing.T) {
	local.TestCase{
		Dir:     "path_arg",
		Command: yamlfmtWithArgs("x.yaml"),
		Update:  *updateFlag,
	}.Run(t)
}

func TestIncludeDocumentStart(t *testing.T) {
	local.TestCase{
		Dir:     "include_document_start",
		Command: yamlfmtWithArgs("-formatter include_document_start=true x.yaml"),
		Update:  *updateFlag,
	}.Run(t)
}

func yamlfmtWithArgs(args string) string {
	return fmt.Sprintf("%s %s", yamlfmtBin, args)
}
