package internal

import (
	"encoding/base64"

	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		// If we are not done
		if m.focused != INSTALL && m.focused != DONE {
			// Setup keybindings
			switch msg.String() {
			case "ctrl+c", "q":
				// Exit when user presses Q or CTRL+C
				return m, tea.Quit

			case "enter":
				if m.width != 0 {
					// Process the selected item, if the return value is true then exit alt screen
					if !m.processSelection() {
						return m, tea.ExitAltScreen
					}
				}
			case "ctrl+z", "backspace":
				// Go backwards in the model
				if m.focused > 0 && m.focused != DONE {
					m.focused--
					return m, nil
				} else {
					// If we are at the beginning, just exit
					return m, tea.Quit
				}
			}
		} else {
			// If we are done then handle keybindings a bit differently
			// Setup keybindings for authDialog
			switch msg.String() {
			case "ctrl+z":
				// Since we have no QuickEmu support, skip the usb controller configuration
				m.focused = VIDEO

			case "ctrl+c":
				// Exit when user presses CTRL+C
				return m, tea.Quit

			case "enter":
				// If we are on the INSTALL dialog
				if m.focused == INSTALL {
					// Write to logger
					logger.Printf("Getting authentication token by elevating with sudo once")

					// Elevate with sudo
					command.Elevate(
						base64.StdEncoding.EncodeToString(
							[]byte(
								m.authDialog.Value(),
							),
						),
					)

					// Write to logger
					logger.Printf("Attempting to free hash from memory")

					// Blank the password field
					m.authDialog.SetValue("")

					// Start installation and send the password to the command
					m.installOutput = m.install()

					// Move to the DONE dialog
					m.focused++

					// Exit the alt screen as the output on the done dialog needs to be scrollable
					//return m, tea.ExitAltScreen

				} else {
					// Quit the application if we are on a different view
					return m, tea.Quit
				}
			}

			// Issue an UI update
			m.authDialog, cmd = m.authDialog.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		if m.width == 0 {
			// Initialize the static lists and make sure the content
			// does not extend past the screen
			m.initLists(msg.Width, msg.Height)

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
					titleStyle = titleStyle.Width(m.width - 4)
					choiceStyle = choiceStyle.Width(m.width)
				}
			}
		}
	}

	// Run another update loop
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}
