package tuimode

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/configs"
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
			PaddingLeft(2)
	choiceStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedChoiceStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))
	dialogStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Width(78)
)

// Make a status type
type status int

// List item struct
type item struct {
	title, desc string
}

// Functions needed for item struct
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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

// Main Model
type model struct {
	fetched    []bool
	lists      []list.Model
	gpu_group  string
	vbios_path string
	loaded     bool
	focused    status
	width      int
	height     int
}

// Consts used to navigate the main model
const (
	INTRO status = iota
	GPUS
	GPU_GROUP
	VBIOS
	VIDEO
	USB
	USB_GROUP
	DONE
)

func (m *model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 10)
	choiceList := list.New([]list.Item{}, choiceDelegate{}, 0, 7)

	// Disable features we wont need
	defaultList.SetShowTitle(false)
	defaultList.SetFilteringEnabled(false)
	defaultList.SetSize(width, height)
	choiceList.SetShowTitle(false)
	choiceList.SetFilteringEnabled(false)

	// Add height and width to our model so we can use it later
	m.width = width
	m.height = height

	m.lists = []list.Model{
		choiceList,
		defaultList,
		defaultList,
		choiceList,
		choiceList,
		defaultList,
		defaultList,
		choiceList,
	}
	m.fetched = []bool{
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
	}
	m.focused = INTRO

	// Init INTRO choices
	items := []list.Item{
		item{title: "CONTINUE"},
	}
	m.lists[INTRO].SetHeight(5)
	m.lists[INTRO].SetItems(items)

	// Init GPU list
	//m.lists[GPUS].Title = "Select a GPU to check the IOMMU groups of"
	items = StringList2ListItem(GetIOMMU("-g", "-F", "name,device_id,optional_revision"))
	m.lists[GPUS].SetItems(items)
	m.fetched[GPUS] = true

	m.lists[GPU_GROUP].Title = ""
	m.lists[GPU_GROUP].SetItems(items)

	// Init USB Controller list
	items = StringList2ListItem(GetIOMMU("-u", "-F", "name,device_id,optional_revision"))
	m.lists[USB].SetItems(items)
	m.fetched[USB] = true

	m.lists[USB_GROUP].Title = ""
	m.lists[USB_GROUP].SetItems(items)

	// Init VBIOS choices
	items = []list.Item{
		item{title: "OK"},
	}
	m.lists[VBIOS].SetItems(items)

	// Init VIDEO disable choises
	items = []list.Item{
		item{title: "YES"},
		item{title: "NO"},
	}
	m.lists[VIDEO].SetItems(items)

	// Init VIDEO disable choises
	items = []list.Item{
		item{title: "FINISH"},
	}
	m.lists[DONE].SetItems(items)
}

func (m model) Init() tea.Cmd {
	return nil
}

// This function processes the enter event
func (m *model) processSelection() {
	switch m.focused {
	case GPUS:
		configs.InitConfigs()

		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Add the gpu group to our model
		m.gpu_group = iommu_group

		items := StringList2ListItem(GetIOMMU("-gr", "-i", m.gpu_group, "-F", "name,device_id,optional_revision"))
		m.lists[GPU_GROUP].SetItems(items)

		// Adjust height to correct for a bigger title
		m.lists[GPU_GROUP].SetSize(m.width, m.height-1)

		// Change focus to next index
		m.focused++

	case GPU_GROUP:
		// Generate the VBIOS dumper script once the user has selected a GPU
		GenerateVBIOSDumper(*m)
		m.focused++

	case USB:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		items := StringList2ListItem(GetIOMMU("-ur", "-i", iommu_group, "-F", "name,device_id,optional_revision"))

		m.lists[USB_GROUP].SetItems(items)

		// Adjust height to correct for a bigger title
		m.lists[USB_GROUP].SetSize(m.width, m.height-1)

		// Change focus to next index
		m.focused++

	case USB_GROUP:
		m.focused++

	case VBIOS:
		m.focused++

	case VIDEO:
		m.focused++

	case INTRO:
		m.focused++

	case DONE:
		os.Exit(0)
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
			//h, v := docStyle.GetFrameSize()

			// Initialize the static lists and make sure the content
			// does not extend past the screen
			m.initLists(msg.Width-2, msg.Height-2)

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
		case INTRO:
			title = dialogStyle.Render(
				fmt.Sprint(
					"Welcome to QuickPassthrough!\n",
					"\n",
					"This script is meant to make it easier to setup GPU passthrough for Qemu systems.\n",
					"However due to the complexity of GPU passthrough, this script assumes you know how to do (or have done) the following.\n\n",
					"* You have already enabled IOMMU, VT-d, SVM and/or AMD-v\n  inside your UEFI/BIOS advanced settings.\n",
					"* Know how to edit your bootloader\n",
					"* Have a bootloader timeout of at least 3 seconds to access the menu\n",
					"* Enable & Configure kernel modules\n",
					"* Have a backup/snapshot of your system in case the script causes your\n  system to be unbootable\n\n",
					"By continuing you accept that I am not liable if your system\n",
					"becomes unbootable, as you will be asked to verify the files generated",
				),
			)
		case GPUS:
			title = titleStyle.Render(
				"Select a GPU to check the IOMMU groups of",
			)

		case GPU_GROUP:
			title = titleStyle.Render(
				fmt.Sprint(
					"Press ENTER/RETURN to set up all these devices for passthrough.\n",
					"This list should only contain items related to your GPU.",
				),
			)

		case USB:
			title = titleStyle.Render(
				"[OPTIONAL]: Select a USB Controller to check the IOMMU groups of",
			)

		case USB_GROUP:
			title = titleStyle.Render(
				fmt.Sprint(
					"Press ENTER/RETURN to set up all these devices for passthrough.\n",
					"This list should only contain the USB controller you want to use.",
				),
			)

		case VBIOS:
			// Get the program directory
			exe, _ := os.Executable()
			scriptdir := filepath.Dir(exe)

			// If we are using go run use the working directory instead
			if strings.Contains(scriptdir, "/tmp/go-build") {
				scriptdir, _ = os.Getwd()
			}

			text := dialogStyle.Render(
				fmt.Sprint(
					"Based on your GPU selection, a vbios extraction script has been generated for your convenience.\n",
					"Passing a VBIOS rom to the card used for passthrough is required for some cards, but not all.\n",
					"Some cards also requires you to patch your VBIOS romfile, check online if this is neccessary for your card!\n",
					"The VBIOS will be read from:\n",
					"%s\n\n",
					"The script to extract the vbios has to be run as sudo and without a displaymanager running for proper dumping!\n",
					"\n",
					"You can run the script with:\n",
					"%s/utils/dump_vbios.sh",
				),
			)

			title = fmt.Sprintf(text, m.vbios_path, scriptdir)

		case VIDEO:
			title = dialogStyle.Render(
				fmt.Sprint(
					"Disabling video output in Linux for the card you want to use in a VM\n",
					"will make it easier to successfully do the passthrough without issues.\n",
					"\n",
					"Do you want to force disable video output in linux on this card?\n",
				),
			)

		case DONE:
			title = dialogStyle.Render(
				fmt.Sprint(
					"The configuration files have been generated and are\n",
					"located inside the \"config\" folder\n",
					"\n",
					"* The \"cmdline\" file contains kernel arguments that your bootloader needs\n",
					"* The \"quickemu\" folder contains files that might be\n  useable for quickemu in the future\n",
					"* The files inside the \"etc\" folder must be copied to your system.\n  NOTE: Verify that these files are correctly formated/edited!\n",
					"\n",
					"A script file named \"install.sh\" has been generated, run it to copy the files to your system and make a backup of your old files.",
				),
			)
		}
		//return listStyle.SetString(fmt.Sprintf("%s\n\n", title)).Render(m.lists[m.focused].View())
		return lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("%s\n%s\n", title, listStyle.Render(m.lists[m.focused].View())))
	} else {
		return "Loading..."
	}
}

func NewModel() *model {
	// Create a blank model and return it
	return &model{}
}

// This is where we build everything
func App() {
	// Make a blank model to keep our state in
	m := NewModel()

	// Start the program with the model
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	errorcheck.ErrorCheck(err, "Failed to initialize UI")
}

func GetIOMMU(args ...string) []string {
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
	var items []string
	output, _ := io.ReadAll(&stdout)

	// Parse the output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		// Write the objects into the list
		items = append(items, scanner.Text())
	}

	// Return our list of items
	return items
}

func StringList2ListItem(stringList []string) []list.Item {
	// Make the []list.Item struct
	items := []list.Item{}

	// Parse the output line by line
	for _, v := range stringList {
		// Get the current line and split by :
		objects := strings.Split(v, ": ")
		// Write the objects into the list
		items = append(items, item{title: objects[1], desc: objects[0]})
	}

	// Return our list of items
	return items
}

func GenerateVBIOSDumper(m model) {
	// Get the vbios path
	m.vbios_path = GetIOMMU("-g", "-i", m.gpu_group, "--rom")[0]

	// Get the config directories
	config := configs.GetConfigPaths()

	// Get the program directory
	exe, _ := os.Executable()
	scriptdir := filepath.Dir(exe)

	// If we are using go run use the working directory instead
	if strings.Contains(scriptdir, "/tmp/go-build") {
		scriptdir, _ = os.Getwd()
	}

	vbios_script_template := fmt.Sprint(
		"#!/bin/bash\n",
		"# THIS FILE IS AUTO GENERATED!\n",
		"# IF YOU HAVE CHANGED GPU, PLEASE RE-RUN QUICKPASSTHROUGH!\n",
		"echo 1 | sudo tee %s\n",
		"sudo bash -c \"cat %s\" > %s/%s/vfio_card.rom\n",
		"echo 0 | sudo tee %s\n",
	)

	vbios_script := fmt.Sprintf(
		vbios_script_template,
		m.vbios_path,
		m.vbios_path,
		scriptdir,
		config.QUICKEMU,
		m.vbios_path,
	)

	scriptfile, err := os.Create("utils/dump_vbios.sh")
	errorcheck.ErrorCheck(err, "Cannot create file \"utils/dump_vbios.sh\"")
	defer scriptfile.Close()

	// Make the script executable
	scriptfile.Chmod(0775)
	errorcheck.ErrorCheck(err, "Could not change permissions of \"utils/dump_vbios.sh\"")

	// Write the script
	scriptfile.WriteString(vbios_script)
}
