package menu

import "github.com/nexidian/gocliselect"

// Make a YesNo menu
func Ok(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("OK", "next")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
