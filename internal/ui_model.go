package internal

import (
	"fmt"
	"os/user"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

// Main Model
type model struct {
	fetched    []bool
	lists      []list.Model
	gpu_group  string
	vbios_path string
	focused    status
	offsetx    []int
	offsety    []int
	width      int
	height     int
	authDialog textinput.Model
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
	INSTALL
	UNKNOWN_BOOTLOADER
	DONE
)

func NewModel() *model {
	// Get the username
	user, err := user.Current()
	errorcheck.ErrorCheck(err, "Error getting username")
	username := user.Username

	// Create the auth input and focus it
	authInput := textinput.New()
	authInput.EchoMode = textinput.EchoPassword
	authInput.Prompt = fmt.Sprintf("\n[sudo] password for %s: ", username)
	authInput.Focus()

	// Create a blank model and return it
	return &model{
		authDialog: authInput,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 10)
	choiceList := list.New([]list.Item{}, choiceDelegate{}, 0, 7)

	// Disable features we wont need
	defaultList.SetShowTitle(false)
	defaultList.SetFilteringEnabled(false)
	defaultList.SetSize(m.width, m.height)
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
		choiceList,
		choiceList,
	}

	// Configure offsets for sizing
	m.offsetx = []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	m.offsety = []int{
		18, 2, 3, 13, 5, 2, 3, 12, 0, 0,
	}

	// Update the styles with the correct width
	dialogStyle = dialogStyle.Width(m.width)
	listStyle = listStyle.Width(m.width)
	titleStyle = titleStyle.Width(m.width - 4)
	choiceStyle = choiceStyle.Width(m.width)

	// Make m.fetched and set all values to FALSE
	m.fetched = []bool{}
	for range m.lists {
		m.fetched = append(m.fetched, false)
	}

	// Set INTRO to the focused view
	m.focused = INTRO

	// Init INTRO choices
	items := []list.Item{
		item{title: "CONTINUE"},
	}
	//m.lists[INTRO].SetHeight(5)
	m.lists[INTRO].SetItems(items)
	m.lists[INTRO].SetSize(m.width-m.offsetx[INTRO], m.height-m.offsety[INTRO])

	// Init GPU list
	items = iommuList2ListItem(getIOMMU("-g", "-F", "vendor:,prod_name,optional_revision:,device_id"))
	m.lists[GPUS].SetItems(items)
	m.lists[GPUS].SetSize(m.width-m.offsetx[GPUS], m.height-m.offsety[GPUS])
	m.fetched[GPUS] = true

	// Setup the initial GPU_GROUP list
	// The content in this list is generated from the selected choice from the GPU view
	m.lists[GPU_GROUP].SetSize(m.width-m.offsetx[GPU_GROUP], m.height-m.offsety[GPU_GROUP])

	// Init USB Controller list
	items = iommuList2ListItem(getIOMMU("-u", "-F", "vendor:,prod_name,optional_revision:,device_id"))
	m.lists[USB].SetItems(items)
	m.lists[USB].SetSize(m.width-m.offsetx[USB], m.height-m.offsety[USB])
	m.fetched[USB] = true

	// Setup the initial USB_GROUP list
	// The content in this list is generated from the selected choice from the USB view
	m.lists[USB_GROUP].SetSize(m.width-m.offsetx[USB_GROUP], m.height-m.offsety[USB_GROUP])

	// Init VBIOS choices
	items = []list.Item{
		item{title: "OK"},
	}
	m.lists[VBIOS].SetItems(items)
	m.lists[VBIOS].SetSize(m.width-m.offsetx[VBIOS], m.height-m.offsety[VBIOS])

	// Init VIDEO disable choises
	items = []list.Item{
		item{title: "YES"},
		item{title: "NO"},
	}
	m.lists[VIDEO].SetItems(items)
	m.lists[VIDEO].SetSize(m.width-m.offsetx[VIDEO], m.height-m.offsety[VIDEO])

	// Init DONE choises
	items = []list.Item{
		item{title: "FINISH"},
	}
	m.lists[DONE].SetItems(items)
	m.lists[DONE].SetSize(m.width-m.offsetx[DONE], m.height-m.offsety[DONE])

	// Init DONE choises
	items = []list.Item{
		item{title: "FINISH"},
	}
	m.lists[UNKNOWN_BOOTLOADER].SetItems(items)
	m.lists[UNKNOWN_BOOTLOADER].SetSize(m.width-m.offsetx[UNKNOWN_BOOTLOADER], m.height-m.offsety[UNKNOWN_BOOTLOADER])
}
