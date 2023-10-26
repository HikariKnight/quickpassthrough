package pages

import (
	"fmt"

	lsiommu "github.com/HikariKnight/quickpassthrough/internal/lsiommu"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
	"github.com/gookit/color"
)

func SelectGPU() {
	// Clear the screen
	command.Clear()

	// Get the users GPUs
	gpus := lsiommu.GetIOMMU("-g", "-F", "vendor:,prod_name,optional_revision:,device_id")

	// Generate a list of choices based on the GPUs and get the users selection
	choice := menu.GenIOMMUMenu("Select a GPU to view the IOMMU groups of", gpus)

	// View the selected GPU
	ViewGPU(choice)
}

func ViewGPU(id string, ext ...int) {
	// Clear the screen
	command.Clear()

	// Set mode to relative
	mode := "-r"

	// Set mode to relative extended
	if len(ext) > 0 {
		mode = "-rr"
	}

	// Get the IOMMU listings for GPUs
	group := lsiommu.GetIOMMU("-g", mode, "-i", id, "-F", "vendor:,prod_name,optional_revision:,device_id")

	// Write a title
	color.Bold.Println("This list should only show devices related to your GPU")

	// Print all the gpus
	for _, v := range group {
		fmt.Println(v)
	}

	// Add a new line for tidyness
	fmt.Println("")

	// Make an empty string
	var choice string

	// Change choices depending on if we have done an extended search or not
	if len(ext) > 0 {
		choice = menu.YesNo("Use this GPU (any extra devices listed may or may not be linked to it) for passthrough?")
	} else {
		choice = menu.YesNoEXT("Use this GPU (and related devices) for passthrough?")
	}

	// Parse the choice
	switch choice {
	case "ext":
		// Run an extended relative search
		ViewGPU(id, 1)

	case "n":
		// Go back to selecting a gpu
		SelectGPU()
	}
}
