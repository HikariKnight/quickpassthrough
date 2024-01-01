package menu

import (
	"github.com/gookit/color"
	"github.com/nexidian/gocliselect"
)

// Make a Next menu
func Next(msg string) string {
	// Make the menu
	menu := gocliselect.NewMenu(msg)
	menu.AddItem(color.Bold.Sprint("Next"), "next")
	menu.AddItem(color.Bold.Sprint("Go Back"), "back")

	// Display the menu
	choice := menu.Display()

	// Return the value selected
	return choice
}
