package menu

import (
	"github.com/gookit/color"
	"github.com/nexidian/gocliselect"
)

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

func YesNoBack(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("Yes", "y")
	menu.AddItem("No", "n")
	menu.AddItem(color.Bold.Sprint("Go Back"), "back")

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
	menu.AddItem("ADVANCED: View with extended related search by vendor ID, devices listed might not be related", "ext")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}

// Make a YesNo menu
func YesNoManual(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("Yes", "y")
	menu.AddItem("No", "n")
	menu.AddItem("Manual Entry", "manual")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
