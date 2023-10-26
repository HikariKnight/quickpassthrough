package menu

import "github.com/nexidian/gocliselect"

// Make a Next menu
func Next(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("Next", "next")
	menu.AddItem("Go Back", "back")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
