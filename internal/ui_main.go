package internal

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/pages"
)

// This is where we build everything
func Tui() {
	// Log all errors to a new logfile (super useful feature of BubbleTea!)
	_ = os.Rename("quickpassthrough_debug.log", "quickpassthrough_debug_old.log")
	logfile, err := tea.LogToFile("quickpassthrough_debug.log", "")
	common.ErrorCheck(err, "Error creating log file")
	defer logfile.Close()

	// New WIP Tui
	pages.Welcome()
}
