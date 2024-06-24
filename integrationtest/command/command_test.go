// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestDry(t *testing.T) {
	TestCase{
		Dir:     "dry",
		Command: yamlfmtWithArgs("-dry ."),
		Update:  *updateFlag,
	}.Run(t)
}

func TestDryQuiet(t *testing.T) {
	TestCase{
		Dir:     "dry_quiet",
		Command: yamlfmtWithArgs("-dry -quiet ."),
		Update:  *updateFlag,
	}.Run(t)
}

func TestPrintConfFlags(t *testing.T) {
	TestCase{
		Dir:     "print_conf_flags",
		Command: yamlfmtWithArgs("-print_conf -continue_on_error=true -formatter retain_line_breaks=true"),
		Update:  *updateFlag,
	}.Run(t)
}

func TestPrintConfFile(t *testing.T) {
	TestCase{
		Dir:     "print_conf_file",
		Command: yamlfmtWithArgs("-print_conf"),
		Update:  *updateFlag,
	}.Run(t)
}
func TestPrintConfFlagsAndFile(t *testing.T) {
	TestCase{
		Dir:     "print_conf_flags_and_file",
		Command: yamlfmtWithArgs("-print_conf -continue_on_error=true -formatter retain_line_breaks=true"),
		Update:  *updateFlag,
	}.Run(t)
}

func TestMultilineStringBug(t *testing.T) {
	TestCase{
		Dir:     "multiline_string_bug",
		Command: yamlfmtWithArgs("-formatter trim_trailing_whitespace=true ."),
		Update:  *updateFlag,
	}.Run(t)
}
