package kyaml

import "github.com/google/yamlfmt"

type KYAMLFormatterFactory struct{}

func (f *KYAMLFormatterFactory) Type() string {
	return KYAMLFormatterType
}

func (f *KYAMLFormatterFactory) NewFormatter(configData map[string]interface{}) (yamlfmt.Formatter, error) {
	return &KYAMLFormatter{}, nil
}
