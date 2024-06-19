package pages

import (
	"strings"
	"testing"
)

func TestFinalizeNotice(t *testing.T) {
	msg := "\n%s\nprinting the finalize notice for manual review, this test should always pass.\n%s\n\n"
	divider := strings.Repeat("-", len(msg)-12)
	t.Logf(msg, divider, divider)
	t.Log("\n\nWith isRoot == true:\n\n")

	finalizeNotice(true)

	println("\n\n")

	t.Log("\n\nWith isRoot == false:\n\n")

	finalizeNotice(false)

	println("\n\n")
}
