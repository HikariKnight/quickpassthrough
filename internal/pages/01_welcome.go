package pages

import (
	"fmt"
	"os"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
	"github.com/gookit/color"
)

// Welcome page
func Welcome() {
	// Clear screen
	command.Clear()

	// Write title
	title := color.New(color.BgHiBlue, color.White, color.Bold)
	title.Println("Welcome to Quickpassthrough!")

	// Write welcome message
	color.Print(
		"This script is meant to make it easier to setup GPU passthrough for\n",
		"Qemu based systems. WITH DIFFERENT 2 GPUS ON THE HOST SYSTEM\n",
		"However due to the complexity of GPU passthrough\n",
		"This script assumes you know how to do (and have done) the following.\n\n",
		"* You have already enabled IOMMU, VT-d, SVM and/or AMD-v\n  inside your UEFI/BIOS advanced settings.\n",
		"* Know how to edit your bootloader\n",
		"* Have a bootloader timeout of at least 3 seconds to access the menu\n",
		"* Enable & Configure kernel modules\n",
		"* Have a backup/snapshot of your system in case the script causes your\n  system to be unbootable\n\n",
		"By continuing you accept that I am not liable if your system\n",
		"becomes unbootable, as you will be asked to verify the files generated\n\n",
		"You can press ESC to exit the program at any time.\n\n",
	)

	// Make user accept responsibility
	choice := menu.YesNo("Are you sure you want to continue?")

	// If yes, go to next page
	if choice == "y" {
		configs.InitConfigs()
		config := configs.GetConfig()
		SelectGPU(config)
	} else {
		fmt.Println("")
		os.Exit(0)
	}
}
