package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fakeSudo = `#!/bin/sh
"$@" -qptest`

const fakeUtil = `#!/bin/sh
echo "$@"
if [ "$4" = "-qptest" ]; then exit 0; else exit 1; fi`

func setupExecTestEnv(t *testing.T) (string, string) {
	t.Helper()

	tmpDir := t.TempDir()
	fakeSudoPath := filepath.Join(tmpDir, "sudo")
	fakeUtilPath := filepath.Join(tmpDir, "util")

	if err := os.WriteFile(fakeSudoPath, []byte(fakeSudo), 0755); err != nil {
		t.Fatalf("failed to write fake sudo stub: %s", err.Error())
	}
	if err := os.WriteFile(fakeUtilPath, []byte(fakeUtil), 0755); err != nil {
		t.Fatalf("failed to write fake util stub: %s", err.Error())
	}
	t.Setenv("PATH", tmpDir+":"+os.Getenv("PATH"))

	return fakeSudoPath, fakeUtilPath
}

func TestExecAndLogSudo(t *testing.T) {
	_, fakeUtilPath := setupExecTestEnv(t)

	args := []string{"i am a string with spaces", "i came to ruin parsers and chew bubble gum", "and I'm all out of bubblegum."}

	t.Run("is_not_root", func(t *testing.T) {
		if err := ExecAndLogSudo(false, false, "util", args...); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("is_root", func(t *testing.T) {
		newFakeUtil := strings.Replace(fakeUtil, "exit 1", "exit 0", 1)
		newFakeUtil = strings.Replace(newFakeUtil, "exit 0", "exit 1", 1)
		if err := os.WriteFile(fakeUtilPath, []byte(newFakeUtil), 0755); err != nil {
			t.Fatalf("failed to overwrite fake util with modified stub: %s", err.Error())
		}
		if err := ExecAndLogSudo(false, false, "util", args...); err == nil {
			t.Errorf("expected error when using modified util with sudo, got nil")
		}

		if err := ExecAndLogSudo(true, true, "util", args...); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

}
