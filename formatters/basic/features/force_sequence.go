// Copyright 2025 Google LLC
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

package features

import "github.com/google/yamlfmt/pkg/yaml"

type SequenceStyle string

const (
	SequenceStyleBlock SequenceStyle = "block"
	SequenceStyleFlow  SequenceStyle = "flow"
)

func FeatureForceSequenceStyle(style SequenceStyle) YAMLFeatureFunc {
	var styleVal yaml.Style
	if style == SequenceStyleFlow {
		styleVal = yaml.FlowStyle
	}
	var forceStyle YAMLFeatureFunc
	forceStyle = func(n yaml.Node) error {
		var err error
		for _, c := range n.Content {
			if c.Kind == yaml.SequenceNode {
				c.Style = styleVal
			}
			err = forceStyle(*c)
		}
		return err
	}
	return forceStyle
}
