// Package gitlab generates GitLab Code Quality reports.
package gitlab

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/google/yamlfmt"
)

// CodeQuality represents a single code quality finding.
//
// Documentation: https://docs.gitlab.com/ee/ci/testing/code_quality.html#code-quality-report-format
type CodeQuality struct {
	Description string
	Name        string
	Fingerprint string
	Severity    Severity
	Path        string
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
		Path:        diff.Path,
	}, true
}

// Marshals a CodeQuality object into JSON.
func (cq CodeQuality) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(wrap(cq))
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return data, nil
}

// UnmarshalJSON unmarshals JSON into a CodeQuality object.
func (cq *CodeQuality) UnmarshalJSON(data []byte) error {
	var ext codeQuality
	if err := json.Unmarshal(data, &ext); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	*cq = unwrap(ext)

	return nil
}

// codeQuality is the external representation of CodeQuality.
// It is needed to add custom JSON marshaling and unmarshaling logic.
type codeQuality struct {
	Description string   `json:"description,omitempty"`
	Name        string   `json:"check_name,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	Severity    Severity `json:"severity,omitempty"`
	Location    location `json:"location,omitempty"`
}

type location struct {
	Path string `json:"path,omitempty"`
}

func wrap(cq CodeQuality) codeQuality {
	return codeQuality{
		Description: cq.Description,
		Name:        cq.Name,
		Fingerprint: cq.Fingerprint,
		Severity:    cq.Severity,
		Location: location{
			Path: cq.Path,
		},
	}
}

func unwrap(cq codeQuality) CodeQuality {
	return CodeQuality{
		Description: cq.Description,
		Name:        cq.Name,
		Fingerprint: cq.Fingerprint,
		Severity:    cq.Severity,
		Path:        cq.Location.Path,
	}
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
