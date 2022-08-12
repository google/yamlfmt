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

package errors

import (
	"fmt"
)

type FormatFileErrors struct {
	formatErrors map[string]error
}

func NewFormatFileErrors() *FormatFileErrors {
	return &FormatFileErrors{
		formatErrors: map[string]error{},
	}
}

func (e *FormatFileErrors) Error() string {
	errStr := "encountered the following formatting errors:"
	for file, err := range e.formatErrors {
		errStr += fmt.Sprintf("%s:%v,", file, err)
	}
	return errStr
}

func (e *FormatFileErrors) Add(file string, err error) {
	e.formatErrors[file] = err
}

func (e *FormatFileErrors) Count() int {
	return len(e.formatErrors)
}

type LintFileErrors struct {
	lintErrors map[string]error
}

func NewLintFileErrors() *LintFileErrors {
	return &LintFileErrors{
		lintErrors: map[string]error{},
	}
}

func (e *LintFileErrors) Error() string {
	errStr := "encountered the following linting errors:\n"
	for file, err := range e.lintErrors {
		errStr += fmt.Sprintf("%s:\n%v\n,", file, err)
	}
	return errStr
}

func (e *LintFileErrors) Add(file string, err error) {
	e.lintErrors[file] = err
}

func (e *LintFileErrors) Count() int {
	return len(e.lintErrors)
}
