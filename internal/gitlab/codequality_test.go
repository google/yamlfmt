package gitlab_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/internal/gitlab"
)

func TestCodeQuality(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		diff            yamlfmt.FileDiff
		wantOK          bool
		wantFingerprint string
	}{
		{
			name: "no diff",
			diff: yamlfmt.FileDiff{
				Path: "testcase/no_diff.yaml",
				Diff: &yamlfmt.FormatDiff{
					Original:  "a: b",
					Formatted: "a: b",
				},
			},
			wantOK: false,
		},
		{
			name: "with diff",
			diff: yamlfmt.FileDiff{
				Path: "testcase/with_diff.yaml",
				Diff: &yamlfmt.FormatDiff{
					Original:  "a:   b",
					Formatted: "a: b",
				},
			},
			wantOK: true,
			// SHA256 of diff.Diff.Original
			wantFingerprint: "05088f1c296b4fd999a1efe48e4addd0f962a8569afbacc84c44630d47f09330",
		},
	}

	for _, tc := range cases {
		// copy tc to avoid capturing an aliased loop variable in a Goroutine.
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, gotOK := gitlab.NewCodeQuality(tc.diff)
			if gotOK != tc.wantOK {
				t.Fatalf("NewCodeQuality() = (%#v, %v), want (*, %v)", got, gotOK, tc.wantOK)
			}
			if !gotOK {
				return
			}

			if tc.wantFingerprint != "" && tc.wantFingerprint != got.Fingerprint {
				t.Fatalf("NewCodeQuality().Fingerprint = %q, want %q", got.Fingerprint, tc.wantFingerprint)
			}

			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}

			var gotUnmarshal gitlab.CodeQuality
			if err := json.Unmarshal(data, &gotUnmarshal); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, gotUnmarshal); diff != "" {
				t.Errorf("json.Marshal() and json.Unmarshal() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
