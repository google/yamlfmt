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
	"os"
	"path/filepath"

	"github.com/google/yamlfmt/command"
)

var (
	flagLint *bool = flag.Bool("lint", false, `Check if there are any differences between
source yaml and formatted yaml.`)
	flagDry *bool = flag.Bool("dry", false, `Perform a dry run; show the output of a formatting
operation without performing it.`)
	flagIn   *bool   = flag.Bool("in", false, "Format yaml read from stdin and output to stdout")
	flagConf *string = flag.String("conf", "", "Read yamlfmt config from this path")
)

func getOperation() command.Operation {
	if *flagIn || isStdinArg() {
		return command.OperationStdin
	}
	if *flagLint {
		return command.OperationLint
	}
	if *flagDry {
		return command.OperationDry
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

func getConfigPath() (string, error) {
	configPath := *flagConf
	if configPath == "" {
		configPath = defaultConfigName
	}

	if filepath.IsAbs(configPath) {
		return configPath, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, configPath), nil
}

func configureHelp() {
	flag.Usage = printHelpMessage
}

func printHelpMessage() {
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
