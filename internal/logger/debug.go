package logger

import (
	"fmt"

	"github.com/google/yamlfmt/internal/collections"
)

type DebugCode int

const (
	DebugCodeAny DebugCode = iota
	DebugCodeConfig
	DebugCodePaths
)

var (
	supportedDebugCodes = map[string][]DebugCode{
		"config": {DebugCodeConfig},
		"paths":  {DebugCodePaths},
		"all":    {DebugCodeConfig, DebugCodePaths},
	}
	activeDebugCodes = collections.Set[DebugCode]{}
)

func ActivateDebugCode(code string) {
	if debugCodes, ok := supportedDebugCodes[code]; ok {
		activeDebugCodes.Add(debugCodes...)
	}
}

// Debug prints a message if the given debug code is active.
func Debug(code DebugCode, msg string, args ...any) {
	if activeDebugCodes.Contains(code) {
		fmt.Printf("[DEBUG]: %s\n", fmt.Sprintf(msg, args...))
	}
}
