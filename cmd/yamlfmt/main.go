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
	"os"
	"path"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/internal/config"
)

var lint *bool = flag.Bool("lint", false, `Check if there are any differences between
source yaml and formatted yaml`)

const defaultConfigName = ".yamlfmt"

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	registry := getFullRegistry()

	var formatter yamlfmt.Formatter
	configData, err := readDefaultConfig()
	if err != nil {
		factory, err := registry.GetFactory(basic.BasicFormatterType)
		if err != nil {
			return err
		}
		formatter = factory.NewDefault()
		if err != nil {
			return err
		}
	} else {
		fType, ok := configData["type"].(string)
		if !ok {
			fType = basic.BasicFormatterType
		}
		factory, err := registry.GetFactory(fType)
		if err != nil {
			return err
		}
		formatter, err = factory.NewWithConfig(configData)
		if err != nil {
			return err
		}
	}

	if *lint {
		return formatter.LintAllFiles()
	}
	return formatter.FormatAllFiles()
}

func readDefaultConfig() (map[string]interface{}, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := path.Join(wd, defaultConfigName)
	configData, err := config.ReadFullConfigFromPath(path)
	if err != nil {
		return nil, err
	}
	return configData, nil
}
