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
	os.Remove("debug.log")
	logfile, err := tea.LogToFile("debug.log", "")
	errorcheck.ErrorCheck(err, "Error creating log file")
	defer logfile.Close()

	// New WIP Tui
	pages.Welcome()
	/*
	   // Make a blank model to keep our state in
	   m := NewModel()

	   // Start the program with the model
	   p := tea.NewProgram(m, tea.WithAltScreen())
	   _, err = p.Run()
	   errorcheck.ErrorCheck(err, "Failed to initialize UI")
	*/
}
