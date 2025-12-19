package kyaml

import (
	"bytes"

	"sigs.k8s.io/yaml/kyaml"
)

const KYAMLFormatterType string = "kyaml"

type KYAMLFormatter struct{}

func (f *KYAMLFormatter) Type() string {
	return KYAMLFormatterType
}

func (f *KYAMLFormatter) Format(input []byte) ([]byte, error) {
	encoder := &kyaml.Encoder{}
	in := bytes.NewReader(input)
	out := bytes.Buffer{}
	if err := encoder.FromYAML(in, &out); err != nil {
		return input, err
	}
	return out.Bytes(), nil
}

func (f *KYAMLFormatter) ConfigMap() (map[string]any, error) {
	return map[string]any{"type": KYAMLFormatterType}, nil
}
