package pages

import (
	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/gookit/color"
)

func prepModules(config *configs.Config) {
	// If we have files for modprobe
	if fileio.FileExist(config.Path.MODPROBE) {
		// Configure modprobe
		configs.Set_Modprobe(config.Gpu_IDs)
	}

	// If we have a folder for dracut
	if fileio.FileExist(config.Path.DRACUT) {
		// Configure dracut
		configs.Set_Dracut()
	}

	// If we have a mkinitcpio.conf file
	if fileio.FileExist(config.Path.MKINITCPIO) {
		configs.Set_Mkinitcpio()
	}

	// Configure grub2 here as we can make the config without sudo
	if config.Bootloader == "grub2" {
		// Write to logger
		logger.Printf("Configuring grub2 manually")
		configs.Configure_Grub2()
	}

	// Finalize changes
	finalize(config)
}

func finalize(config *configs.Config) {
	// Clear the screen
	command.Clear()

	// Write a title
	title := color.New(color.BgHiBlue, color.White, color.Bold)
	title.Println("Finalizing configuration")

	color.Print(
		"The configuration files have been generated and are\n",
		"located inside the \"config\" folder\n",
		"\n",
		"* The \"kernel_args\" file contains kernel arguments that your bootloader needs\n",
		//"* The \"quickemu\" folder contains files that might be\n  useable for quickemu in the future\n",
		"* The files inside the \"etc\" folder must be copied to your system.\n",
		"  NOTE: Verify that these files are correctly formated/edited!\n",
		"* Once all files have been copied, you need to update your bootloader and rebuild\n",
		"  your initramfs using the tools to do so by your system.\n",
		"\n",
		"This program can do this for you, however the program will have to\n",
		"type your password to sudo using STDIN, to avoid using STDIN press CTRL+C\n",
		"and copy the files, update your bootloader and rebuild your initramfs manually.\n",
		"If you want to go back and change something, press CTRL+Z\n",
		"\nNOTE: A backup of the original files from the first run can be found in the backup folder\n",
	)

}
