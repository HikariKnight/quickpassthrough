package menu

import "github.com/nexidian/gocliselect"

// Make an OK menu
func Ok(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("OK", "next")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}

// Make an OK & Go Back menu
func OkBack(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem("OK", "next")
	menu.AddItem("Go Back", "back")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
