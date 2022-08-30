package hotfix_test

import (
	"testing"

	"github.com/google/yamlfmt/formatters/basic"
	"github.com/google/yamlfmt/internal/hotfix"
)

func TestParseEmoji(t *testing.T) {
	testCases := []struct {
		name        string
		yamlStr     string
		expectedStr string
	}{
		{
			name:        "parses emoji",
			yamlStr:     "a: 😂\n",
			expectedStr: "a: \"😂\"\n",
		},
		{
			name:        "parses multiple emoji",
			yamlStr:     "a: 😼 👑\n",
			expectedStr: "a: \"😼 👑\"\n",
		},
	}

	f := &basic.BasicFormatter{Config: basic.DefaultConfig()}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formattedBefore, err := f.Format([]byte(tc.yamlStr))
			if err != nil {
				t.Fatalf("yaml failed to parse: %v", err)
			}
			formattedAfter := hotfix.ParseUnicodePoints(formattedBefore)
			formattedStr := string(formattedAfter)
			if formattedStr != tc.expectedStr {
				t.Fatalf("parsed string does not match: \nexpected: %s\ngot: %s", tc.expectedStr, string(formattedStr))
			}
		})
	}
}
