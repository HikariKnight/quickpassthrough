package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyToSystem(t *testing.T) {
	if err := os.Mkdir("testdir", 0755); err != nil {
		t.Fatal(err)
	}
	tFilePath := filepath.Join("testdir", "testfile")
	if err := os.WriteFile(tFilePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll("testdir"); err != nil {
			t.Fatal(err)
		}
	})
	isRoot := os.Getuid() == 0
	switch isRoot {
	case true:
		t.Run("TestCopyToSystem_AsRoot", func(t *testing.T) {
			CopyToSystem(true, tFilePath, "/etc/testfile")
		})
	default:
		t.Run("TestCopyToSystem_AsUser", func(t *testing.T) {
			CopyToSystem(false, tFilePath, "/etc/testfile")
		})
	}
}
