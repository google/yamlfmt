package features

import (
	"errors"
	"fmt"

	"github.com/google/yamlfmt/pkg/yaml"
)

type QuoteStyle string

const (
	SingleQuoteStyle QuoteStyle = "single"
	DoubleQuoteStyle QuoteStyle = "double"
)

var ErrUnrecognizedQuoteStyle = errors.New("unrecognized quote style")

func FeatureForceQuoteStyle(style QuoteStyle) (YAMLFeatureFunc, error) {
	var fromStyle, toStyle yaml.Style
	switch style {
	case SingleQuoteStyle:
		fromStyle = yaml.DoubleQuotedStyle
		toStyle = yaml.SingleQuotedStyle
	case DoubleQuoteStyle:
		fromStyle = yaml.SingleQuotedStyle
		toStyle = yaml.DoubleQuotedStyle
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnrecognizedQuoteStyle, style)
	}
	var forceStyle YAMLFeatureFunc
	forceStyle = func(n yaml.Node) error {
		var err error
		for _, c := range n.Content {
			if c.Style == fromStyle {
				c.Style = toStyle
			}
			err = forceStyle(*c)
		}
		return err
	}
	return forceStyle, nil
}
