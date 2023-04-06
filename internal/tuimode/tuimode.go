package tuimode

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

// This is where we handle all the events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// This is where we change what is shown on the screen
func (m model) View() string {
	return "test\n"
}

// This is where we build everything
func App() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
