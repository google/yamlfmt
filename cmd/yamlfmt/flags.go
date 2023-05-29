// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/google/yamlfmt/command"
)

var (
	flagLint *bool = flag.Bool("lint", false, `Check if there are any differences between
source yaml and formatted yaml.`)
	flagDry *bool = flag.Bool("dry", false, `Perform a dry run; show the output of a formatting
operation without performing it.`)
	flagIn              *bool   = flag.Bool("in", false, "Format yaml read from stdin and output to stdout")
	flagVersion         *bool   = flag.Bool("version", false, "Print yamlfmt version")
	flagConf            *string = flag.String("conf", "", "Read yamlfmt config from this path")
	flagDoublestar      *bool   = flag.Bool("dstar", false, "Use doublestar globs for include and exclude")
	flagQuiet           *bool   = flag.Bool("quiet", false, "Print minimal output to stdout")
	flagContinueOnError *bool   = flag.Bool("continue_on_error", false, "Continue to format files that didn't fail instead of exiting with code 1.")
	flagExclude                 = arrayFlag{}
	flagFormatter               = arrayFlag{}
	flagExtensions              = arrayFlag{}
)

func bindArrayFlags() {
	flag.Var(&flagExclude, "exclude", "Paths to exclude in the chosen format (standard or doublestar)")
	flag.Var(&flagFormatter, "formatter", "Config value overrides to pass to the formatter")
	flag.Var(&flagExtensions, "extensions", "File extensions to use for standard path collection")
}

type arrayFlag []string

// Implements flag.Value
func (a *arrayFlag) String() string {
	return strings.Join(*a, " ")
}

func (a *arrayFlag) Set(value string) error {
	values := []string{value}
	if strings.Contains(value, ",") {
		values = strings.Split(value, ",")
	}
	*a = append(*a, values...)
	return nil
}

func configureHelp() {
	flag.Usage = func() {
		fmt.Println(`yamlfmt is a simple command line tool for formatting yaml files.

	Arguments:

	Glob paths to yaml files
			Send any number of paths to yaml files specified in doublestar glob format (see: https://github.com/bmatcuk/doublestar). 
			Any flags must be specified before the paths.

	- or /dev/stdin
			Passing in a single - or /dev/stdin will read the yaml from stdin and output the formatted result to stdout
		
	Flags:`)
		flag.PrintDefaults()
	}
}

func getOperationFromFlag() command.Operation {
	if *flagIn || isStdinArg() {
		return command.OperationStdin
	}
	if *flagLint {
		return command.OperationLint
	}
	if *flagDry {
		return command.OperationDry
	}
	if *flagVersion {
		return command.OperationVersion
	}
	return command.OperationFormat
}

func isStdinArg() bool {
	if len(flag.Args()) != 1 {
		return false
	}
	arg := flag.Args()[0]
	return arg == "-" || arg == "/dev/stdin"
}
