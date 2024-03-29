package menu

import (
	"regexp"

	"github.com/gookit/color"
	"github.com/nexidian/gocliselect"
)

func GenIOMMUMenu(msg string, choices []string, none_option ...int) string {
	// Make a regex to get the iommu group
	iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)

	// Make the menu
	menu := gocliselect.NewMenu(msg)

	// For each choice passed
	for _, choice := range choices {
		// Get the iommu group
		iommuGroup := iommu_group_regex.FindString(choice)

		// Add the choice with shortened vendor name and the iommu group as the return value
		menu.AddItem(choice, iommuGroup)
	}

	// If none_option is higher than 0
	if len(none_option) > 0 {
		// Add a skip option
		menu.AddItem(color.Bold.Sprint("Skip/None"), "skip")
	}

	// Add a go back option
	menu.AddItem(color.Bold.Sprint("Go Back"), "back")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
