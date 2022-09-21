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

import (
	"fmt"
)

type Formatter interface {
	Type() string
	Format(yamlContent []byte) ([]byte, error)
}

type Factory interface {
	Type() string
	NewFormatter(config map[string]interface{}) (Formatter, error)
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

type LineBreakStyle string

const (
	LineBreakStyleLF   LineBreakStyle = "lf"
	LineBreakStyleCRLF LineBreakStyle = "crlf"
)

type UnsupportedLineBreakError struct {
	style LineBreakStyle
}

func (e UnsupportedLineBreakError) Error() string {
	return fmt.Sprintf("unsupported line break style %s, see package documentation for supported styles", e.style)
}

func (s LineBreakStyle) Separator() (string, error) {
	switch s {
	case LineBreakStyleLF:
		return "\n", nil
	case LineBreakStyleCRLF:
		return "\r\n", nil
	}
	return "", UnsupportedLineBreakError{style: s}
}

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

type FeatureApplyError struct {
	err         error
	featureName string
	mode        FeatureApplyMode
}

func (e *FeatureApplyError) Error() string {
	action := "Before"
	if e.mode == FeatureApplyAfter {
		action = "After"
	}
	return fmt.Sprintf("Feature %s %sAction failed with error: %v", e.featureName, action, e.err)
}

func (e *FeatureApplyError) Unwrap() error {
	return e.err
}

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

		if err != nil {
			return nil, &FeatureApplyError{
				err:         err,
				featureName: feature.Name,
				mode:        mode,
			}
		}
	}
	return result, nil
}
