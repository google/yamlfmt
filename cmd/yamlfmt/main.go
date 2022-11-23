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
	"log"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/command"
	"github.com/google/yamlfmt/formatters/basic"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	configureHelp()
	flag.Parse()

	operation := getOperationFromFlag()

	configData := map[string]any{}
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	if configPath != "" {
		configData, err = readConfig(configPath)
		if err != nil {
			return err
		}
	}

	if len(flag.Args()) > 0 {
		configData["include"] = flag.Args()
	}

	return command.RunCommand(operation, getFullRegistry(), configData)
}

func getFullRegistry() *yamlfmt.Registry {
	return yamlfmt.NewFormatterRegistry(&basic.BasicFormatterFactory{})
}
