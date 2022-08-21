package command

import (
	"fmt"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/engine"
	"github.com/google/yamlfmt/formatters/basic"
	"github.com/mitchellh/mapstructure"
)

type Operation int

const (
	OperationFormat Operation = iota
	OperationLint
	OperationDry
)

type commandConfig struct {
	Include         []string               `mapstructure:"include"`
	Exclude         []string               `mapstructure:"exclude"`
	FormatterConfig map[string]interface{} `mapstructure:"formatter"`
}

func RunCommand(
	operation Operation,
	registry *yamlfmt.Registry,
	configData map[string]interface{},
) error {
	var config commandConfig
	err := mapstructure.Decode(configData, &config)
	if err != nil {
		return err
	}
	if len(config.Include) == 0 {
		config.Include = []string{"**/*.{yaml,yml}"}
	}

	var formatter yamlfmt.Formatter
	if config.FormatterConfig == nil {
		factory, err := registry.GetDefaultFactory()
		if err != nil {
			return err
		}
		formatter = factory.NewDefault()
	} else {
		fType, ok := config.FormatterConfig["type"].(string)
		if !ok {
			fType = basic.BasicFormatterType
		}
		factory, err := registry.GetFactory(fType)
		if err != nil {
			return err
		}
		formatter, err = factory.NewWithConfig(config.FormatterConfig)
		if err != nil {
			return err
		}
	}

	engine := &engine.Engine{
		Include:   config.Include,
		Exclude:   config.Exclude,
		Formatter: formatter,
	}

	switch operation {
	case OperationFormat:
		err := engine.FormatAllFiles()
		if err != nil {
			return err
		}
	case OperationLint:
		err := engine.LintAllFiles()
		if err != nil {
			return err
		}
	case OperationDry:
		out, err := engine.DryRunAllFiles()
		if err != nil {
			return err
		}
		fmt.Println(out)
	}

	return nil
}
