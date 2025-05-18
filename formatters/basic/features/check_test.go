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

package features_test

import (
	"strings"
	"testing"

	"github.com/google/yamlfmt/formatters/basic/features"
	"github.com/google/yamlfmt/pkg/yaml"
)

func TestCheck(t *testing.T) {
	for _, c := range []struct {
		desc    string
		in      string
		wantErr bool
	}{{
		desc:    "no anchors",
		in:      "foo: bar",
		wantErr: false,
	}, {
		desc: "anchor",
		// From https://support.atlassian.com/bitbucket-cloud/docs/yaml-anchors/
		in: `
definitions:
  steps:
    - step: &build-test
        name: Build and test
        script:
          - mvn package
        artifacts:
          - target/**

pipelines:
  branches:
    develop:
      - step: *build-test
    main:
      - step: *build-test
`,
		wantErr: true,
	}} {
		t.Run(c.desc, func(t *testing.T) {
			var docNode yaml.Node
			if err := yaml.NewDecoder(strings.NewReader(c.in)).Decode(&docNode); err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if err := features.Check(docNode); (err != nil) != c.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
