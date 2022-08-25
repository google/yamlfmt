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
	"errors"
	"flag"
	"log"
	"os"
	"path"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/command"
	"github.com/google/yamlfmt/formatters/basic"
	"gopkg.in/yaml.v3"
)

var (
	flagLint *bool = flag.Bool("lint", false, `Check if there are any differences between
source yaml and formatted yaml.`)
	flagDry *bool = flag.Bool("dry", false, `Perform a dry run; show the output of a formatting
operation without performing it.`)
	flagConf *string = flag.String("conf", "", "Read yamlfmt config from this path")
)

const defaultConfigName = ".yamlfmt"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()

	operation := getOperation()

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	configData, err := readConfig(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			configData = map[string]interface{}{}
		} else {
			return err
		}
	}

	if len(flag.Args()) > 0 {
		configData["include"] = flag.Args()
	}

	return command.RunCommand(operation, getFullRegistry(), configData)
}

func getConfigPath() (string, error) {
	configPath := *flagConf
	if configPath == "" {
		configPath = defaultConfigName
	}

	if path.IsAbs(configPath) {
		return configPath, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(wd, configPath), nil
}

func readConfig(path string) (map[string]interface{}, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var configData map[string]interface{}
	err = yaml.Unmarshal(yamlBytes, &configData)
	if err != nil {
		return nil, err
	}
	return configData, nil
}

func getOperation() command.Operation {
	if len(flag.Args()) == 1 && isStdin(flag.Args()[0]) {
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

func isStdin(arg string) bool {
	return arg == "-" || arg == "/dev/stdin"
}

func getFullRegistry() *yamlfmt.Registry {
	return yamlfmt.NewFormatterRegistry(&basic.BasicFormatterFactory{})
}
