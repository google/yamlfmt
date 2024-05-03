// Copyright 2024 Google LLC
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

package yamlfmt_test

import (
	"errors"
	"testing"

	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/collections"
)

type errorChecker func(t *testing.T, errs collections.Errors)

func checkErrNil(t *testing.T, errs collections.Errors) {
	combinedErr := errs.Combine()
	if combinedErr != nil {
		t.Fatalf("expected error to be nil, got: %v", combinedErr)
	}
}

func TestReadMetadata(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected collections.Set[yamlfmt.Metadata]
		errCheck errorChecker
	}{
		{
			name:     "contains no metadata",
			content:  "",
			expected: collections.Set[yamlfmt.Metadata]{},
			errCheck: checkErrNil,
		},
		{
			name:    "has ignore metadata",
			content: "# !yamlfmt!:ignore\na: 1",
			expected: collections.Set[yamlfmt.Metadata]{
				{Type: yamlfmt.MetadataIgnore, LineNum: 1}: {},
			},
			errCheck: checkErrNil,
		},
		{
			name:     "has bad metadata",
			content:  "# !yamlfmt!fjghgh",
			expected: collections.Set[yamlfmt.Metadata]{},
			errCheck: func(t *testing.T, errs collections.Errors) {
				if len(errs) != 1 {
					t.Fatalf("expected 1 error, got %d:\n%v", len(errs), errs.Combine())
				}
				if errors.Unwrap(errs[0]) != yamlfmt.ErrMalformedMetadata {
					t.Fatalf("expected ErrMalformedMetadata, got: %v", errs[0])
				}
			},
		},
		{
			name:     "has unrecognized metadata type",
			content:  "# !yamlfmt!:lulsorandom",
			expected: collections.Set[yamlfmt.Metadata]{},
			errCheck: func(t *testing.T, errs collections.Errors) {
				if len(errs) != 1 {
					t.Fatalf("expected 1 error, got %d:\n%v", len(errs), errs.Combine())
				}
				if errors.Unwrap(errs[0]) != yamlfmt.ErrUnrecognizedMetadata {
					t.Fatalf("expected ErrUnrecognizedMetadata, got: %v", errs[0])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			md, err := yamlfmt.ReadMetadata([]byte(tc.content), "test.yaml")
			if !md.Equals(tc.expected) {
				t.Fatalf("Mismatched metadata:\nexpected: %v\ngot: %v", tc.expected, md)
			}
			t.Logf("got error: %v", err)
			tc.errCheck(t, err)
		})
	}
}
