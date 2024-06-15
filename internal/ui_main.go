package internal

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"os"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/pages"
	tea "github.com/charmbracelet/bubbletea"
)

// This is where we build everything
func Tui() {
	// Log all errors to a new logfile (super useful feature of BubbleTea!)
	_ = os.Rename("quickpassthrough_debug.log", "quickpassthrough_debug_old.log")
	logfile, err := tea.LogToFile("quickpassthrough_debug.log", "")
	errorcheck.ErrorCheck(err, "Error creating log file")
	defer logfile.Close()

	// New WIP Tui
	pages.Welcome()
}
