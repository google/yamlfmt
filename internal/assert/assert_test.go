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

package assert_test

import (
	"fmt"
	"testing"

	"github.com/google/yamlfmt/internal/assert"
)

type tMock struct {
	logs   []string
	failed bool
	err    error
}

func newTMock() *tMock {
	return &tMock{
		logs: []string{},
	}
}

func (t *tMock) Helper() {}

func (t *tMock) Fatal(...any) {
	t.failed = true
}

func (t *tMock) Fatalf(msg string, args ...any) {
	t.logs = append(t.logs, fmt.Sprintf(msg, args...))
	t.Fatal()
}

func (t *tMock) Errorf(msg string, args ...any) {
	t.failed = true
	t.err = fmt.Errorf(msg, args...)
}

func TestAssertFail(t *testing.T) {
	testInstance := newTMock()
	failMsg := "expected %d to equal %d"
	a := 1
	b := 2
	assert.Assert(testInstance, a == b, failMsg, a, b)
	if !testInstance.failed {
		t.Fatalf("Assert failed. %v", *testInstance)
	}
	if len(testInstance.logs) != 1 {
		t.Fatalf("Found %d logs. %v", len(testInstance.logs), testInstance.logs)
	}
	expectedFailLog := fmt.Sprintf(failMsg, a, b)
	if testInstance.logs[0] != expectedFailLog {
		t.Fatalf(
			"Failure log didn't match.\nexpected: %s\ngot: %s",
			expectedFailLog,
			testInstance.logs[0],
		)
	}
}

func TestEqualFail(t *testing.T) {
	testInstance := newTMock()
	failMsg := "expected %v to equal %v"
	expected := 1
	got := 2
	assert.EqualMsg(testInstance, expected, got, failMsg)
	if len(testInstance.logs) != 1 {
		t.Fatalf("Found %d logs. %v", len(testInstance.logs), testInstance.logs)
	}
	expectedFailLog := fmt.Sprintf(failMsg, expected, got)
	if testInstance.logs[0] != expectedFailLog {
		t.Fatalf(
			"Failure log didn't match.\nexpected: %s\ngot: %s",
			expectedFailLog,
			testInstance.logs[0],
		)
	}
}

func TestDereferenceEqualErr(t *testing.T) {
	testInstance := newTMock()
	expected := &struct{}{}
	errMsg := "nil pointer %v %v"
	assert.DereferenceEqualMsg(testInstance, expected, nil, errMsg, "does not matter")
	if testInstance.err == nil {
		t.Fatalf("DereferenceEqual should have failed")
	}
	expectedErr := fmt.Errorf(errMsg, expected, nil)
	if testInstance.err.Error() != expectedErr.Error() {
		t.Fatalf(
			"Errors didn't match.\nexpected: %s\ngot: %s",
			expectedErr,
			testInstance.err,
		)
	}
}

func TestDerefenceEqualFail(t *testing.T) {
	testInstance := newTMock()
	type x struct {
		num int
	}
	failMsg := "%v not equal %v"
	expected := &x{num: 1}
	got := &x{num: 2}
	assert.DereferenceEqualMsg(testInstance, expected, got, "does not matter", failMsg)
	if len(testInstance.logs) != 1 {
		t.Fatalf("Found %d logs. %v", len(testInstance.logs), testInstance.logs)
	}
	expectedFailLog := fmt.Sprintf(failMsg, *expected, *got)
	if testInstance.logs[0] != expectedFailLog {
		t.Fatalf(
			"Failure log didn't match.\nexpected: %s\ngot: %s",
			expectedFailLog,
			testInstance.logs[0],
		)
	}
}

func TestDereferenceEqualPass(t *testing.T) {
	testInstance := newTMock()
	type x struct {
		num int
	}
	expected := &x{num: 1}
	got := &x{num: 1}
	assert.DereferenceEqualMsg(testInstance, expected, got, "doesn't matter", "doesn't matter")
	if testInstance.failed {
		t.Fatalf("test failed when it should have passed")
	}
	if len(testInstance.logs) != 0 {
		t.Fatalf("test instance had logs when it shouldn't: %v", testInstance.logs)
	}
}

func TestSliceEqualFailDiffSize(t *testing.T) {
	testInstance := newTMock()
	failSizeMsg := "%v and %v"
	expected := []int{1, 2, 3, 4}
	got := []int{1, 2, 3}
	assert.SliceEqualMsg(testInstance, expected, got, failSizeMsg, "something else")
	if len(testInstance.logs) != 1 {
		t.Fatalf("Found %d logs. %v", len(testInstance.logs), testInstance.logs)
	}
	expectedFailLog := fmt.Sprintf(failSizeMsg, len(expected), len(got))
	if testInstance.logs[0] != expectedFailLog {
		t.Fatalf(
			"Failure log didn't match.\nexpected: %s\ngot: %s",
			expectedFailLog,
			testInstance.logs[0],
		)
	}
}

func TestSliceEqualMismatch(t *testing.T) {
	testInstance := newTMock()
	failMismatchMsg := "at index %v: %v and %v"
	expected := []int{1, 2, 4}
	got := []int{1, 2, 3}
	assert.SliceEqualMsg(testInstance, expected, got, "something else", failMismatchMsg)
	if len(testInstance.logs) != 1 {
		t.Fatalf("Found %d logs. %v", len(testInstance.logs), testInstance.logs)
	}
	expectedFailLog := fmt.Sprintf(failMismatchMsg, 2, expected[2], got[2])
	if testInstance.logs[0] != expectedFailLog {
		t.Fatalf(
			"Failure log didn't match.\nexpected: %s\ngot: %s",
			expectedFailLog,
			testInstance.logs[0],
		)
	}
}
