#!/usr/bin/env bash

# Copyright 2024 GitLab, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu -o pipefail

if [[ $# -ne 1 ]]; then
	echo "Usage: $0 <testname>" >&2
	exit 1
fi

TESTNAME="$1"
shift

declare -r BEFORE_DIR="integrationtest/command/testdata/${TESTNAME:?}/before"
declare -r AFTER_DIR="integrationtest/command/testdata/${TESTNAME:?}/after"
declare -r OUTPUT_DIR="integrationtest/command/testdata/${TESTNAME:?}/stdout"
declare -r GO_TEST_FILE="integrationtest/command/command_test.go"

if [[ ! -d "${BEFORE_DIR}" ]]; then
	mkdir -p "${BEFORE_DIR}"
	cat >"${BEFORE_DIR}/x.yaml" <<EOF
TODO: add test input
EOF
fi

if [[ ! -d "${AFTER_DIR}" ]]; then
	mkdir -p "${AFTER_DIR}"
	cat >"${AFTER_DIR}/x.yaml" <<EOF
TODO: add test input
EOF
fi

mkdir -p "${OUTPUT_DIR}"
touch "${OUTPUT_DIR}/stdout.txt"
touch "${OUTPUT_DIR}/stderr.txt"

# Converty "my_test_case" to "MyTestCase".
if sed --version 2>/dev/null | grep -q GNU; then
	# \U (convert to upper case) is a GNU extension.
	GO_TESTNAME="$(sed -r -e 's/(^|_)(\w)/\U\2/g' <<<"${TESTNAME}")"
else
	# Fall back to Perl.
	GO_TESTNAME="$(perl -pe 's/(^|_)(\w)/\U\2/g' <<<"${TESTNAME}")"
fi

if ! grep -q -F "Test${GO_TESTNAME:?}" "${GO_TEST_FILE}"; then
	cat >>"${GO_TEST_FILE}" <<EOF

func Test${GO_TESTNAME:?}(t *testing.T) {
	TestCase{
		Dir:     "${TESTNAME:?}",
		// TODO: Change arguments to match your test case.
		Command: yamlfmtWithArgs("."),
		Update:  *updateFlag,
	}.Run(t)
}
EOF
	gofmt -w "${GO_TEST_FILE}"
fi

cat <<EOF
Next steps:

* Write test input files to '${BEFORE_DIR}'
* Write expected output to '${OUTPUT_DIR}/stdout.txt'
* Optional: write expected error output to '${OUTPUT_DIR}/stderr.txt'
* Add yamlfmt arguments to the 'Test${GO_TESTNAME}' test in '${GO_TEST_FILE}'
* Run 'make integrationtest_update' and review the files in '${AFTER_DIR}'
EOF
