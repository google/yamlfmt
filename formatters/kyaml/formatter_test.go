// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kyaml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/yamlfmt/internal/assert"
)

func TestKYAMLFormatter(t *testing.T) {
	// This might need to be changed to a proper
	// struct table if it tests anything other than
	// simple output equality.
	testCases := []string{
		"basic_case",
	}
	for _, testName := range testCases {
		t.Run(testName, func(t *testing.T) {
			f := &KYAMLFormatter{}
			testdataPath := filepath.Join("testdata", testName)
			before, err := os.ReadFile(filepath.Join(testdataPath, "before.yaml"))
			assert.NilErr(t, err)
			after, err := os.ReadFile(filepath.Join(testdataPath, "after.yaml"))
			assert.NilErr(t, err)
			result, err := f.Format(before)
			assert.NilErr(t, err)
			assert.Equal(t, string(result), string(after))
		})
	}
}
