package pages

import (
	"fmt"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
)

func disableVideo() {
	// Clear the screen
	command.Clear()

	// Get our config struct
	//config := configs.GetConfig()

	fmt.Print(
		"Disabling video output in Linux for the card you want to use in a VM\n",
		"will make it easier to successfully do the passthrough without issues.\n",
		"\n",
	)

	// Make the yesno menu
	choice := menu.YesNo("Do you want to force disable video output in linux on this card?")

	if choice == "Yes" {
		// Add disable VFIO video to the config
		configs.DisableVFIOVideo(1)
	} else {
		// Do not disable VFIO Video
		configs.DisableVFIOVideo(0)
	}

}
