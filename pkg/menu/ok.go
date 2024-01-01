package menu

import (
	"github.com/gookit/color"
	"github.com/nexidian/gocliselect"
)

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
	menu.AddItem(color.Bold.Sprint("OK"), "next")
	menu.AddItem(color.Bold.Sprint("Go Back"), "back")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
