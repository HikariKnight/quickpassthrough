package pages

import (
	"fmt"
	"os"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
	"github.com/gookit/color"
)

func disableVideo(config *configs.Config) {
	// Clear the screen
	command.Clear()

	// Write a title
	title := color.New(color.BgHiBlue, color.White, color.Bold)
	title.Println("Do you want to disable video output on the VFIO card in Linux?")

	fmt.Print(
		"Disabling video output in Linux for the card you want to use in a VM\n",
		"will make it easier to successfully do the passthrough without issues.\n",
		"\n",
	)

	// Make the yesno menu
	choice := menu.YesNoBack("Do you want to force disable video output in linux on this card?")

	switch choice {
	case "y":
		// Add disable VFIO video to the config
		configs.DisableVFIOVideo(1)
		//selectUSB(config)
		prepModules(config)

	case "n":
		// Do not disable VFIO Video
		configs.DisableVFIOVideo(0)
		//selectUSB(config)
		prepModules(config)

	case "back":
		genVBIOS_dumper(config)

	case "":
		// If ESC is pressed
		fmt.Println("")
		os.Exit(0)
	}
}
