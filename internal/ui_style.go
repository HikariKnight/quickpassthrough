package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#5F5FD7")).
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(2).PaddingRight(2)
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(241))
	listStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2)
	choiceStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingRight(4)
	selectedChoiceStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))
	dialogStyle = lipgloss.NewStyle().
			PaddingLeft(2)
)

// Choice delegate (for our dialog boxes)
type choiceDelegate struct{}

func (d choiceDelegate) Height() int                               { return 1 }
func (d choiceDelegate) Spacing() int                              { return 0 }
func (d choiceDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d choiceDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.title

	fn := choiceStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedChoiceStyle.Render("| " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
