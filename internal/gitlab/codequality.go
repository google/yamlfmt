// Package gitlab generates GitLab Code Quality reports.
package gitlab

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/yamlfmt"
)

// CodeQuality represents a single code quality finding.
//
// Documentation: https://docs.gitlab.com/ee/ci/testing/code_quality.html#code-quality-report-format
type CodeQuality struct {
	Description string   `json:"description,omitempty"`
	Name        string   `json:"check_name,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	Severity    Severity `json:"severity,omitempty"`
	Location    Location `json:"location,omitempty"`
}

// Location is the location of a Code Quality finding.
type Location struct {
	Path string `json:"path,omitempty"`
}

// NewCodeQuality creates a new CodeQuality object from a yamlfmt.FileDiff.
//
// If the file did not change, i.e. the diff is empty, an empty struct and false is returned.
func NewCodeQuality(diff yamlfmt.FileDiff) (CodeQuality, bool) {
	if !diff.Diff.Changed() {
		return CodeQuality{}, false
	}

	return CodeQuality{
		Description: "Not formatted correctly, run yamlfmt to resolve.",
		Name:        "yamlfmt",
		Fingerprint: fingerprint(diff),
		Severity:    Major,
		Location: Location{
			Path: diff.Path,
		},
	}, true
}

// fingerprint returns a 256-bit SHA256 hash of the original unformatted file.
// This is used to uniquely identify a code quality finding.
func fingerprint(diff yamlfmt.FileDiff) string {
	hash := sha256.New()

	fmt.Fprint(hash, diff.Diff.Original)

	return fmt.Sprintf("%x", hash.Sum(nil)) //nolint:perfsprint
}

// Severity is the severity of a code quality finding.
type Severity string

const (
	Info     Severity = "info"
	Minor    Severity = "minor"
	Major    Severity = "major"
	Critical Severity = "critical"
	Blocker  Severity = "blocker"
)
