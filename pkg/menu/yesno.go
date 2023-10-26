package menu

import "github.com/nexidian/gocliselect"

// Make a YesNo menu
func YesNo(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("Yes", "y")
	menu.AddItem("No", "n")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}

func YesNoEXT(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("Yes", "y")
	menu.AddItem("No", "n")
	menu.AddItem("ADVANCED: View with extended related search by vendor ID, results will be inaccurate", "ext")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
