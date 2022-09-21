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

package basic_test

import (
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/formatters/basic"
)

func TestNewWithConfigRetainsDefaultValues(t *testing.T) {
	testCases := []struct {
		name           string
		configMap      map[string]interface{}
		expectedConfig basic.Config
	}{
		{
			name: "only indent specified",
			configMap: map[string]interface{}{
				"indent": 4,
			},
			expectedConfig: basic.Config{
				Indent:               4,
				IncludeDocumentStart: false,
				LineEnding:           yamlfmt.LineBreakStyleLF,
			},
		},
		{
			name: "only include_document_start specified",
			configMap: map[string]interface{}{
				"include_document_start": true,
			},
			expectedConfig: basic.Config{
				Indent:               2,
				IncludeDocumentStart: true,
				LineEnding:           yamlfmt.LineBreakStyleLF,
			},
		},
		{
			name: "only line_ending style specified",
			configMap: map[string]interface{}{
				"line_ending": "crlf",
			},
			expectedConfig: basic.Config{
				Indent:               2,
				IncludeDocumentStart: false,
				LineEnding:           yamlfmt.LineBreakStyleCRLF,
			},
		},
		{
			name: "all specified",
			configMap: map[string]interface{}{
				"indent":                 4,
				"line_ending":            "crlf",
				"include_document_start": true,
			},
			expectedConfig: basic.Config{
				Indent:               4,
				IncludeDocumentStart: true,
				LineEnding:           yamlfmt.LineBreakStyleCRLF,
			},
		},
	}

	factory := basic.BasicFormatterFactory{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter, err := factory.NewFormatter(tc.configMap)
			if err != nil {
				t.Fatalf("expected factory to create config, got error: %v", err)
			}
			basicFormatter, ok := formatter.(*basic.BasicFormatter)
			if !ok {
				t.Fatal("should have been able to cast to basic formatter")
			}
			if *basicFormatter.Config != tc.expectedConfig {
				t.Fatalf("configs differed:\nexpected: %v\ngot: %v", *basicFormatter.Config, tc.expectedConfig)
			}
		})
	}
}
