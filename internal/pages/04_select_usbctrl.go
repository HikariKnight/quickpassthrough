package pages

import (
	"fmt"
	"os"

	lsiommu "github.com/HikariKnight/quickpassthrough/internal/lsiommu"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
	"github.com/gookit/color"
)

func selectUSB() {
	// Clear the screen
	command.Clear()

	// Get the users GPUs
	usbs := lsiommu.GetIOMMU("-u", "-F", "vendor:,prod_name,optional_revision:,device_id")

	// Generate a list of choices based on the GPUs and get the users selection
	choice := menu.GenIOMMUMenu("Select a USB to view the IOMMU groups of", usbs, 1)

	// Parse the choice
	switch choice {
	case "back":
		SelectGPU()
	case "":
		// If ESC is pressed
		fmt.Println("")
		os.Exit(0)
	default:
		// View the selected GPU
		viewUSB(choice)
	}
}

func viewUSB(id string, ext ...int) {
	// Clear the screen
	command.Clear()

	// Set mode to relative
	mode := "-r"

	// Set mode to relative extended
	if len(ext) > 0 {
		mode = "-rr"
	}

	// Get the IOMMU listings for USB controllers
	group := lsiommu.GetIOMMU("-u", mode, "-i", id, "-F", "vendor:,prod_name,optional_revision:,device_id")

	// Write a title
	color.Bold.Println("This list should only show the USB controller")

	// Print all the usb controllers
	for _, v := range group {
		fmt.Println(v)
	}

	// Add a new line for tidyness
	fmt.Println("")

	// Make an empty string
	var choice string

	// Ask if we shall use the devices for passthrough
	if len(ext) == 0 {
		choice = menu.YesNo("Use all listed devices for passthrough?")
	} else {
		choice = menu.YesNoEXT("Use all listed devices for passthrough?")
	}

	// Parse the choice
	switch choice {
	case "":
		// If ESC is pressed
		fmt.Println("")
		os.Exit(0)
	case "n":
		// Go back to selecting a gpu
		selectUSB()

	case "y":
		// Go to the select a usb controller
	}
}
