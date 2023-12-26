package pages

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"syscall"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
	"github.com/HikariKnight/quickpassthrough/pkg/uname"
	"github.com/gookit/color"
	"golang.org/x/term"
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
		"If you want to go back and change something, choose Back\n",
		"\nNOTE: A backup of the original files from the first run can be found in the backup folder\n",
	)

	// Make a choice of going next or back
	choice := menu.Next("Press Next to continue with sudo using STDIN, ESC to exit or Back to go back.")

	// Parse the choice
	switch choice {
	case "next":
		installPassthrough(config)

	case "back":
		// Go back
		disableVideo(config)
	}

}

func installPassthrough(config *configs.Config) {
	// Get the user data
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Provide a password prompt
	fmt.Printf("[sudo] password for %s: ", user.Username)
	bytep, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	fmt.Print("\n")

	// Elevate with sudo
	command.Elevate(
		base64.StdEncoding.EncodeToString(
			bytep,
		),
	)

	// Make an output string
	var output string

	// Based on the bootloader, setup the configuration
	if config.Bootloader == "kernelstub" {
		// Write to logger
		logger.Printf("Configuring systemd-boot using kernelstub")

		// Configure kernelstub
		output = configs.Set_KernelStub()
		fmt.Printf("%s\n", output)

	} else if config.Bootloader == "grubby" {
		// Write to logger
		logger.Printf("Configuring bootloader using grubby")

		// Configure kernelstub
		output = configs.Set_Grubby()
		fmt.Printf("%s\n", output)

	} else if config.Bootloader == "grub2" {
		// Write to logger
		logger.Printf("Configuring grub2 manually")
		grub_output, _ := configs.Set_Grub2()
		fmt.Printf("%s\n", strings.Join(grub_output, "\n"))

	} else {
		kernel_args := fileio.ReadFile(config.Path.CMDLINE)
		logger.Printf("Unsupported bootloader, please add the below line to your bootloaders kernel arguments\n%s", kernel_args)
	}

	// A lot of linux systems support modprobe along with their own module system
	// So copy the modprobe files if we have them
	modprobeFile := fmt.Sprintf("%s/vfio.conf", config.Path.MODPROBE)
	if fileio.FileExist(modprobeFile) {
		// Copy initramfs-tools module to system
		output = configs.CopyToSystem(modprobeFile, "/etc/modprobe.d/vfio.conf")
		fmt.Printf("%s\n", output)
	}

	// Copy the config files for the system we have
	initramfsFile := fmt.Sprintf("%s/modules", config.Path.INITRAMFS)
	dracutFile := fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT)
	if fileio.FileExist(initramfsFile) {
		// Copy initramfs-tools module to system
		output = configs.CopyToSystem(initramfsFile, "/etc/initramfs-tools/modules")
		fmt.Printf("%s\n", output)

		// Copy the modules file to /etc/modules
		output = configs.CopyToSystem(config.Path.ETCMODULES, "/etc/modules")
		fmt.Printf("%s\n", output)

		// Write to logger
		logger.Printf("Executing: sudo update-initramfs -u")

		// Update initramfs
		fmt.Println("Executed: sudo update-initramfs -u\nSee debug.log for detailed output")
		cmd_out, cmd_err, _ := command.RunErr("sudo", "update-initramfs", "-u")

		cmd_out = append(cmd_out, cmd_err...)

		// Write to logger
		logger.Printf(strings.Join(cmd_out, "\n"))
	} else if fileio.FileExist(dracutFile) {
		// Copy dracut config to /etc/dracut.conf.d/vfio
		output = configs.CopyToSystem(dracutFile, "/etc/dracut.conf.d/vfio")
		fmt.Printf("%s\n", output)

		// Get systeminfo
		sysinfo := uname.New()

		// Write to logger
		logger.Printf("Executing: sudo dracut -f -v --kver %s", sysinfo.Release)

		// Update initramfs
		fmt.Printf("Executed: sudo dracut -f -v --kver %s\nSee debug.log for detailed output", sysinfo.Release)
		cmd_out, cmd_err, _ := command.RunErr("sudo", "dracut", "-f", "-v", "--kver", sysinfo.Release)

		cmd_out = append(cmd_out, cmd_err...)

		// Write to logger
		logger.Printf(strings.Join(cmd_out, "\n"))
	} else if fileio.FileExist(config.Path.MKINITCPIO) {
		// Copy dracut config to /etc/dracut.conf.d/vfio
		output = configs.CopyToSystem(config.Path.MKINITCPIO, "/etc/mkinitcpio.conf")
		fmt.Printf("%s\n", output)

		// Write to logger
		logger.Printf("Executing: sudo mkinitcpio -P")

		// Update initramfs
		fmt.Println("Executed: sudo mkinitcpio -P\nSee debug.log for detailed output")
		cmd_out, cmd_err, _ := command.RunErr("sudo", "mkinitcpio", "-P")

		cmd_out = append(cmd_out, cmd_err...)

		// Write to logger
		logger.Printf(strings.Join(cmd_out, "\n"))
	}
}
