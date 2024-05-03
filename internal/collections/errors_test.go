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

package collections_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/yamlfmt/internal/collections"
)

func TestErrorsCombine(t *testing.T) {
	errs := collections.Errors{
		errors.New("a"),
		nil,
		errors.New("c"),
	}
	err := errs.Combine()
	if err == nil {
		t.Fatal("expected combined err not to be nil")
	}
	for _, errEl := range errs {
		if errEl == nil {
			continue
		}
		if !strings.Contains(err.Error(), errEl.Error()) {
			t.Fatalf("expected combined err to contain %v, got: %v", errEl, err)
		}
	}
}

func TestErrorsCombineEmpty(t *testing.T) {
	errs := collections.Errors{}
	err := errs.Combine()
	if err != nil {
		t.Fatalf("expected combined err to be nil, got: %v", err)
	}
}

func TestErrorsCombineNilElements(t *testing.T) {
	errs := collections.Errors{nil, nil, nil}
	err := errs.Combine()
	if err != nil {
		t.Fatalf("expected combined err to be nil, got: %v", err)
	}
}
