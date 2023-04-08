package internal

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	tea "github.com/charmbracelet/bubbletea"
)

// This is where we build everything
func Tui() {
	// Make a blank model to keep our state in
	m := NewModel()

	// Start the program with the model
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	errorcheck.ErrorCheck(err, "Failed to initialize UI")
}
