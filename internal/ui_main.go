package internal

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Setup keybindings
		switch msg.String() {
		case "ctrl+c", "q":
			// Exit when user presses Q or CTRL+C
			return m, tea.Quit

		case "enter":
			if m.loaded {
				// Process the selected item
				m.processSelection()
			}
		case "ctrl+z", "backspace":
			// Go backwards in the model
			if m.focused > 0 {
				m.focused--
				return m, nil
			} else {
				// If we are at the beginning, just exit
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		if !m.loaded {
			// Initialize the static lists and make sure the content
			// does not extend past the screen
			m.initLists(msg.Width, msg.Height)

			// Set model loaded to true
			m.loaded = true
		} else {
			// Else we are loaded and will update the sizing on the fly
			m.height = msg.Height
			m.width = msg.Width

			// TODO: Find a better way to resize widgets when word wrapping happens
			// BUG: currently breaks the UI rendering if word wrapping happens in some cases...
			views := len(m.lists)
			if msg.Width > 83 {
				for i := 0; i < views; i++ {
					m.lists[i].SetSize(m.width-m.offsetx[i], m.height-m.offsety[i])
					// Update the styles with the correct width
					dialogStyle = dialogStyle.Width(m.width)
					listStyle = listStyle.Width(m.width)
					titleStyle = titleStyle.Width(m.width - 2)
					choiceStyle = choiceStyle.Width(m.width)
				}
			}
		}
	}
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

// This is where we build everything
func Tui() {
	// Make a blank model to keep our state in
	m := NewModel()

	// Start the program with the model
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	errorcheck.ErrorCheck(err, "Failed to initialize UI")
}
