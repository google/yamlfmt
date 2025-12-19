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
