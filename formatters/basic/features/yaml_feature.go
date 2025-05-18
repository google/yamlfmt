package features

import "github.com/google/yamlfmt/pkg/yaml"

// These features will directly use the `yaml.Node` type and
// as such are specific to this formatter.
type YAMLFeatureFunc func(yaml.Node) error
type YAMLFeatureList []YAMLFeatureFunc

func (y YAMLFeatureList) ApplyFeatures(node yaml.Node) error {
	for _, f := range y {
		if err := f(node); err != nil {
			return err
		}
	}
	return nil
}
