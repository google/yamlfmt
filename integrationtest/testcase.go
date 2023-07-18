package integrationtest

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/yamlfmt/internal/assert"
	"github.com/google/yamlfmt/internal/tempfile"
)

const (
	stdoutGoldenFile = "stdout.txt"
)

type tempDirTestCase struct {
	Dir     string
	Command string
	Update  bool
}

func (tc tempDirTestCase) Run(t *testing.T) {
	// I wanna write on the first indent level lol
	t.Run(tc.Dir, tc.run)
}

func (tc tempDirTestCase) run(t *testing.T) {
	// Replicate the "before" directory in the test temp directory.
	tempDir := t.TempDir()
	paths, err := tempfile.ReplicateDirectory(tc.testFolderBeforePath(), tempDir)
	assert.NilErr(t, err)
	err = paths.CreateAll()
	assert.NilErr(t, err)

	// Run the command for the test in the temp directory.
	var stdoutBuf bytes.Buffer
	cmd := tc.command(tempDir, &stdoutBuf)
	err = cmd.Run()
	assert.NilErr(t, err)

	err = tc.goldenStdout(stdoutBuf.Bytes())
	assert.NilErr(t, err)
	err = tc.goldenAfter(tempDir)
	assert.NilErr(t, err)
}

func (tc tempDirTestCase) testFolderBeforePath() string {
	return tc.testdataDirPath() + "/before"
}

func (tc tempDirTestCase) command(wd string, stdoutBuf *bytes.Buffer) *exec.Cmd {
	cmdParts := strings.Split(tc.Command, " ")
	return &exec.Cmd{
		Path:   cmdParts[0], // This is just the path to look up the binary
		Args:   cmdParts,    // Args needs to be an array of everything including the command name
		Stdout: stdoutBuf,
		Dir:    wd,
	}
}

func (tc tempDirTestCase) goldenStdout(stdoutResult []byte) error {
	goldenCtx := tempfile.GoldenCtx{
		Dir:    tc.testFolderStdoutPath(),
		Update: tc.Update,
	}
	return goldenCtx.CompareGoldenFile(stdoutGoldenFile, stdoutResult)
}

func (tc tempDirTestCase) goldenAfter(wd string) error {
	goldenCtx := tempfile.GoldenCtx{
		Dir:    tc.testFolderAfterPath(),
		Update: tc.Update,
	}
	return goldenCtx.CompareDirectory(wd)
}

func (tc tempDirTestCase) testFolderAfterPath() string {
	return filepath.Join(tc.testdataDirPath(), "after")
}

func (tc tempDirTestCase) testFolderStdoutPath() string {
	return filepath.Join(tc.testdataDirPath(), "stdout")
}

func (tc tempDirTestCase) testdataDirPath() string {
	return filepath.Join("testdata/local", tc.Dir)
}
