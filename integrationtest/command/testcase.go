//go:build integration_test

package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/yamlfmt/internal/assert"
	"github.com/google/yamlfmt/internal/tempfile"
)

const (
	stdoutGoldenFile = "stdout.txt"
	stderrGoldenFile = "stderr.txt"
)

type TestCase struct {
	Dir        string
	Command    string
	IsError    bool
	Update     bool
	ShowStdout bool
}

func (tc TestCase) Run(t *testing.T) {
	// I wanna write on the first indent level lol
	t.Run(tc.Dir, tc.run)
}

func (tc TestCase) run(t *testing.T) {
	// Replicate the "before" directory in the test temp directory.
	tempDir := t.TempDir()
	paths, err := tempfile.ReplicateDirectory(tc.testFolderBeforePath(), tempDir)
	assert.NilErr(t, err)
	err = paths.CreateAll()
	assert.NilErr(t, err)

	// Run the command for the test in the temp directory.
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd := tc.command(tempDir, &stdoutBuf, &stderrBuf)
	err = cmd.Run()
	if tc.IsError {
		assert.NotNilErr(t, err)
	} else {
		assert.NilErr(t, err)
	}

	fmt.Printf("stdout: %s\n", stdoutBuf.String())
	fmt.Printf("stderr: %s\n", stderrBuf.String())
	err = tc.goldenStdout(stdoutBuf.Bytes())
	assert.NilErr(t, err)
	err = tc.goldenStderr(stderrBuf.Bytes())
	assert.NilErr(t, err)
	err = tc.goldenAfter(tempDir)
	assert.NilErr(t, err)
}

func (tc TestCase) testFolderBeforePath() string {
	return tc.testdataDirPath() + "/before"
}

func (tc TestCase) command(wd string, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) *exec.Cmd {
	cmdArgs := []string{}
	for _, arg := range strings.Split(tc.Command, " ") {
		// This is to handle potential typos in args with extra spaces.
		if arg != "" {
			cmdArgs = append(cmdArgs, arg)
		}
	}
	return &exec.Cmd{
		Path:   cmdArgs[0], // This is just the path to the command
		Args:   cmdArgs,    // Args needs to be an array of everything including the command
		Stdout: stdoutBuf,
		Stderr: stderrBuf,
		Dir:    wd,
	}
}

func (tc TestCase) goldenStdout(stdoutResult []byte) error {
	if tc.ShowStdout {
		fmt.Printf("Output for test %s:\n%s", tc.Dir, stdoutResult)
		return nil
	}
	goldenCtx := tempfile.GoldenCtx{
		Dir:    tc.testFolderStdoutPath(),
		Update: tc.Update,
	}
	return goldenCtx.CompareGoldenFile(stdoutGoldenFile, stdoutResult)
}

func (tc TestCase) goldenStderr(stderrResult []byte) error {
	if tc.ShowStdout {
		fmt.Printf("Stderr output for test %s:\n%s", tc.Dir, stderrResult)
		return nil
	}
	goldenCtx := tempfile.GoldenCtx{
		Dir:    tc.testFolderStdoutPath(),
		Update: tc.Update,
	}
	return goldenCtx.CompareGoldenFile(stderrGoldenFile, stderrResult)
}

func (tc TestCase) goldenAfter(wd string) error {
	goldenCtx := tempfile.GoldenCtx{
		Dir:    tc.testFolderAfterPath(),
		Update: tc.Update,
	}
	return goldenCtx.CompareDirectory(wd)
}

func (tc TestCase) testFolderAfterPath() string {
	return filepath.Join(tc.testdataDirPath(), "after")
}

func (tc TestCase) testFolderStdoutPath() string {
	return filepath.Join(tc.testdataDirPath(), "stdout")
}

func (tc TestCase) testdataDirPath() string {
	return filepath.Join("testdata/", tc.Dir)
}
