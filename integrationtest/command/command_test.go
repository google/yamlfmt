//go:build integration_test

package command

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	updateFlag *bool = flag.Bool("update", false, "Whether to update the goldens.")
	stdoutFlag *bool = flag.Bool("stdout", false, "Show stdout instead of diffing it.")
	yamlfmtBin string
)

func TestMain(m *testing.M) {
	yamlfmtBinVar := os.Getenv("YAMLFMT_BIN")
	if yamlfmtBinVar == "" {
		fmt.Println("Must provide a YAMLFMT_BIN environment variable.")
		os.Exit(1)
	}
	yamlfmtBin = yamlfmtBinVar
	m.Run()
}

func yamlfmtWithArgs(args string) string {
	return fmt.Sprintf("%s -no_global_conf %s", yamlfmtBin, args)
}

func TestPathArg(t *testing.T) {
	TestCase{
		Dir:     "path_arg",
		Command: yamlfmtWithArgs("x.yaml"),
		Update:  *updateFlag,
	}.Run(t)
}

func TestIncludeDocumentStart(t *testing.T) {
	TestCase{
		Dir:     "include_document_start",
		Command: yamlfmtWithArgs("-formatter include_document_start=true x.yaml"),
		Update:  *updateFlag,
	}.Run(t)
}

func TestGitignore(t *testing.T) {
	TestCase{
		Dir:     "gitignore",
		Command: yamlfmtWithArgs("-gitignore_excludes -gitignore_path .test_gitignore ."),
		Update:  *updateFlag,
	}.Run(t)
}

func TestLint(t *testing.T) {
	TestCase{
		Dir:     "lint",
		Command: yamlfmtWithArgs("-lint ."),
		Update:  *updateFlag,
		IsError: true,
	}.Run(t)
}

func TestLineOutput(t *testing.T) {
	TestCase{
		Dir:     "line_output",
		Command: yamlfmtWithArgs("-lint -output_format line ."),
		Update:  *updateFlag,
		IsError: true,
	}.Run(t)
}
