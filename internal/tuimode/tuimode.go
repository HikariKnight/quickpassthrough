package tuimode

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle   = lipgloss.NewStyle().Margin(2, 2)
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#5F5FD7")).
			Foreground(lipgloss.Color("#FFFFFF"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(241))
	listStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder())
)

type status int

const (
	GPUS status = iota
	GPU_GROUP
	USB
	USB_GROUP
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	fetched []bool
	lists   []list.Model
	loaded  bool
	focused status
	width   int
	height  int
}

func (m *model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height/2)

	// Disable features we wont need
	defaultList.SetShowTitle(false)
	defaultList.SetFilteringEnabled(false)
	defaultList.SetSize(width, height)

	// Add height and width to our model so we can use it later
	m.width = width
	m.height = height

	m.lists = []list.Model{defaultList, defaultList, defaultList, defaultList}
	m.fetched = []bool{false, false, false, false}
	m.focused = GPUS

	// Init GPU list
	//m.lists[GPUS].Title = "Select a GPU to check the IOMMU groups of"
	items := GetIOMMU("-g", "-F", "name,device_id,optional_revision")
	m.lists[GPUS].SetShowTitle(false)
	m.lists[GPUS].SetItems(items)
	m.fetched[GPUS] = true

	m.lists[GPU_GROUP].Title = ""
	m.lists[GPU_GROUP].SetItems(items)

	// Init USB Controller list
	items = GetIOMMU("-u", "-F", "name,device_id,optional_revision")
	m.lists[USB].SetItems(items)
	m.fetched[USB] = true

	m.lists[USB_GROUP].Title = ""
	m.lists[USB_GROUP].SetItems(items)
}

func (m model) Init() tea.Cmd {
	return nil
}

// This function processes the enter event
func (m *model) processSelection() {
	switch m.focused {
	case GPUS:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		items := GetIOMMU("-gr", "-i", iommu_group, "-F", "name,device_id,optional_revision")
		m.lists[GPU_GROUP].SetItems(items)

		// Adjust height to correct for a bigger title
		m.lists[GPU_GROUP].SetSize(m.width, m.height)

		// Change focus to next index
		m.focused++

	case GPU_GROUP:
		// Gets the selected item
		/*selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		items := GetIOMMU("-gr", "-i", iommu_group, "--id")*/
		m.focused++
	case USB:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		items := GetIOMMU("-ur", "-i", iommu_group, "-F", "name,device_id,optional_revision")
		m.lists[USB_GROUP].SetItems(items)

		// Adjust height to correct for a bigger title
		m.lists[USB_GROUP].SetSize(m.width, m.height)

		// Change focus to next index
		m.focused++
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.loaded {
				m.processSelection()
			}
		case "ctrl+z", "backspace":
			if m.focused > 0 {
				m.focused--
				return m, nil
			} else {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		if !m.loaded {
			// Get the terminal frame size
			h, v := docStyle.GetFrameSize()

			// Initialize the static lists and make sure the content
			// does not extend past the screen
			m.initLists(msg.Width-h, msg.Height-v)

			// Set model loaded to true
			m.loaded = true
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.loaded {
		title := ""
		switch m.focused {
		case GPUS:
			title = fmt.Sprintf(
				"\n %s\n",
				titleStyle.Render("Select a GPU to check the IOMMU groups of"),
			)

		case GPU_GROUP:
			title = fmt.Sprintf(
				"\n %s\n %s",
				titleStyle.Render("Press ENTER/RETURN to set up all these devices for passthrough."), titleStyle.Render("This list should only contain items related to your GPU."),
			)

		case USB:
			title = fmt.Sprintf(
				"\n %s\n",
				titleStyle.Render("[OPTIONAL]: Select a USB Controller to check the IOMMU groups of"),
			)

		case USB_GROUP:
			title = fmt.Sprintf(
				"\n %s\n %s",
				titleStyle.Render("Press ENTER/RETURN to set up all these devices for passthrough."), titleStyle.Render("This list should only contain the USB controller you want to use."),
			)
		}
		return lipgloss.JoinVertical(lipgloss.Left, title, listStyle.Render(m.lists[m.focused].View()))
	} else {
		return "Loading..."
	}
}

func NewModel() *model {
	return &model{}
}

// This is where we build everything
func App() {
	m := NewModel()

	// Start the program with the model
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func GetIOMMU(args ...string) []list.Item {
	var stdout, stderr bytes.Buffer

	// Configure the ls-iommu command
	cmd := exec.Command("utils/ls-iommu", args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	// Execute the command
	err := cmd.Run()

	// If ls-iommu returns an error then IOMMU is disabled
	errorcheck.ErrorCheck(err, "IOMMU disabled in either UEFI/BIOS or in bootloader!")

	// Read the output
	items := []list.Item{}
	output, _ := io.ReadAll(&stdout)

	// Parse the output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		// Get the current line and split by :
		objects := strings.Split(scanner.Text(), ": ")
		// Write the objects into the list
		items = append(items, item{title: objects[1], desc: objects[0]})
	}

	// Return our list of items
	return items
}
