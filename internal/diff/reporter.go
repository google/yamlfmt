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

package diff

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type prettyReporter struct {
	path  cmp.Path
	lines []string
}

func (r *prettyReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *prettyReporter) Report(rs cmp.Result) {
	vx, vy := r.path.Last().Values()
	if !rs.Equal() {
		if vx.IsValid() {
			r.lines = append(r.lines, fmt.Sprintf("- %+v", vx))
		}
		if vy.IsValid() {
			r.lines = append(r.lines, fmt.Sprintf("+ %+v", vy))
		}
	} else {
		r.lines = append(r.lines, fmt.Sprintf("  %+v", vx))
	}
}

func (r *prettyReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *prettyReporter) String() string {
	return strings.Join(r.lines, "\n")
}
