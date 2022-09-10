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

package yamlfmt

import "fmt"

type Formatter interface {
	Type() string
	Format(yamlContent []byte) ([]byte, error)
}

type Factory interface {
	Type() string
	NewDefault() Formatter
	NewWithConfig(config map[string]interface{}) (Formatter, error)
}

type Registry struct {
	registry    map[string]Factory
	defaultType string
}

func NewFormatterRegistry(defaultFactory Factory) *Registry {
	return &Registry{
		registry: map[string]Factory{
			defaultFactory.Type(): defaultFactory,
		},
		defaultType: defaultFactory.Type(),
	}
}

func (r *Registry) Add(f Factory) {
	r.registry[f.Type()] = f
}

func (r *Registry) GetFactory(fType string) (Factory, error) {
	factory, ok := r.registry[fType]
	if !ok {
		return nil, fmt.Errorf("no formatter registered with type \"%s\"", fType)
	}
	return factory, nil
}

func (r *Registry) GetDefaultFactory() (Factory, error) {
	factory, ok := r.registry[r.defaultType]
	if !ok {
		return nil, fmt.Errorf("no default formatter registered for type \"%s\"", r.defaultType)
	}
	return factory, nil
}

const (
	LineBreakStyleLF   = "lf"
	LineBreakStyleCRLF = "crlf"
)

type FeatureFunc func([]byte) ([]byte, error)

type Feature struct {
	Name         string
	BeforeAction FeatureFunc
	AfterAction  FeatureFunc
}

type FeatureList []Feature

type FeatureApplyMode string

var (
	FeatureApplyBefore FeatureApplyMode = "before"
	FeatureApplyAfter  FeatureApplyMode = "after"
)

func (fl FeatureList) ApplyFeatures(input []byte, mode FeatureApplyMode) ([]byte, error) {
	// Declare err here so the result variable doesn't get shadowed in the loop
	var err error
	result := make([]byte, len(input))
	copy(result, input)
	for _, feature := range fl {
		if mode == FeatureApplyBefore {
			if feature.BeforeAction != nil {
				result, err = feature.BeforeAction(result)
			}
		} else {
			if feature.AfterAction != nil {
				result, err = feature.AfterAction(result)
			}
		}

		// TODO: use the yamlfmt.Feature error wrap function I'll make
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
