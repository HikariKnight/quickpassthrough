package tuimode

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	lists map[string]list.Model
	list  list.Model
}

func initLists() model {
	m := model{lists: make(map[string]list.Model)}
	return m
}

func (m model) addItems(obj string, items []list.Item) {
	m.lists[obj] = list.New(items, list.NewDefaultDelegate(), 0, 0)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

// This is where we build everything
func App() {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("utils/ls-iommu", "-g", "-F", "name,device_id,optional_revision")
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	errorcheck.ErrorCheck(err)

	items := []list.Item{}
	output, _ := io.ReadAll(&stdout)

	// If we get an error from ls-iommu then print the error and exit
	// As the system might not have IOMMU enabled
	if stderr.String() != "" {
		log.Fatal(stderr.String())
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		objects := strings.Split(scanner.Text(), ": ")
		items = append(items, item{title: objects[1], desc: objects[0]})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a GPU to check the IOMMU groups of"

	m := model{list: l}

	// Start the program with the model
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
