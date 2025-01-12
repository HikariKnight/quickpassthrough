package pages

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"syscall"

	"github.com/gookit/color"
	"golang.org/x/term"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
	"github.com/HikariKnight/quickpassthrough/pkg/uname"
)

func prepModules(config *configs.Config) {
	// If we have files for modprobe
	if exists, _ := fileio.FileExist(config.Path.MODPROBE); exists {
		// Configure modprobe
		configs.Set_Modprobe(config.Gpu_IDs)
	}

	if exists, _ := fileio.FileExist(config.Path.INITRAMFS); exists && config.HasDuplicateDeviceIds {
		// Configure initramfs early binds
		configs.SetInitramfsToolsEarlyBinds(config)
	}

	// If we have a folder for dracut
	if exists, _ := fileio.FileExist(config.Path.DRACUT); exists {
		// Configure dracut
		configs.Set_Dracut(config)
	}

	// If we have a mkinitcpio.conf file
	if exists, _ := fileio.FileExist(config.Path.MKINITCPIO); exists {
		configs.Set_Mkinitcpio(config)
	}

	// Configure grub2 here as we can make the config without sudo
	if config.Bootloader == "grub2" {
		// Write to logger
		logger.Printf("Configuring grub2 manually\n")
		configs.Configure_Grub2()
	}

	// Finalize changes
	finalize(config)
}

func finalizeNotice(isRoot bool) {
	color.Print(`
The configuration files have been generated and are located inside the "config" folder

  * The "kernel_args" file contains kernel arguments that your bootloader needs
  * The "qemu" folder contains files that may be needed for passthrough
  * The files inside the "etc" folder must be copied to your system.

	<red>Verify that these files are correctly formated/edited!</>

Once all files have been copied, the following steps must be taken:

  * bootloader configuration must be updated
  * initramfs must be rebuilt

`)
	switch isRoot {
	case true:
		color.Print("This program can do this for you, if desired.\n")
	default:
		color.Print(`This program can do this for you, however your sudo password is required.
To avoid this:

  * press CTRL+C and perform the steps mentioned above manually.
       OR
  * run ` + os.Args[0] + ` as root.

`)
	}

	color.Print(`
If you want to go back and change something, choose Back.

NOTE: A backup of the original files from the first run can be found in the backup folder
`)
}

func finalize(config *configs.Config) {
	// Clear the screen
	command.Clear()

	// Write a title
	title := color.New(color.BgHiBlue, color.White, color.Bold)
	title.Println("Finalizing configuration")

	config.IsRoot = os.Getuid() == 0

	finalizeNotice(config.IsRoot)

	// Make a choice of going next or back and parse the choice
	switch menu.Next("Press Next to continue with sudo using STDIN, ESC to exit or Back to go back.") {
	case "next":
		installPassthrough(config)
	case "back":
		// Go back
		disableVideo(config)
	}
}

func installPassthrough(config *configs.Config) {
	// Get the user data
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if !config.IsRoot {
		// Provide a password prompt
		fmt.Printf("[sudo] password for %s: ", currentUser.Username)
		bytep, err := term.ReadPassword(syscall.Stdin)
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
	}

	// Make an output string
	var output string

	// Based on the bootloader, setup the configuration
	if config.Bootloader == "kernelstub" {
		// Write to logger
		logger.Printf("Configuring systemd-boot using kernelstub\n")

		// Configure kernelstub
		// callee logs the output and checks for errors
		configs.Set_KernelStub(config.IsRoot)

	} else if config.Bootloader == "grubby" {
		// Write to logger
		logger.Printf("Configuring bootloader using grubby\n")

		// Configure kernelstub
		output = configs.Set_Grubby(config.IsRoot)
		fmt.Printf("%s\n", output)

	} else if config.Bootloader == "grub2" {
		// Write to logger
		logger.Printf("Applying grub2 changes\n")
		_ = configs.Set_Grub2(config.IsRoot) // note: we set config.IsRoot earlier

		// we'll print the output in the [configs.Set_Grub2] method
		// fmt.Printf("%s\n", strings.Join(grub_output, "\n"))

	} else {
		kernel_args := fileio.ReadFile(config.Path.CMDLINE)
		logger.Printf("Unsupported bootloader, please add the below line to your bootloaders kernel arguments\n%s", kernel_args)
	}

	// A lot of linux systems support modprobe along with their own module system
	// So copy the modprobe files if we have them
	modprobeFile := fmt.Sprintf("%s/vfio.conf", config.Path.MODPROBE)

	// lets hope by now we've already handled any permissions issues...
	// TODO: verify that we actually can drop the errors on [fileio.FileExist] call below

	if exists, _ := fileio.FileExist(modprobeFile); exists {
		// Copy initramfs-tools module to system, note that CopyToSystem will log the command and output
		// as well as check for errors
		configs.CopyToSystem(config.IsRoot, modprobeFile, "/etc/modprobe.d/vfio.conf")
	}

	// Copy the config files for the system we have
	initramfsFile := fmt.Sprintf("%s/modules", config.Path.INITRAMFS)
	dracutFile := fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT)

	initramFsExists, initramFsErr := fileio.FileExist(initramfsFile)
	dracutExists, dracutErr := fileio.FileExist(dracutFile)
	mkinitcpioExists, mkinitcpioErr := fileio.FileExist(config.Path.MKINITCPIO)

	for _, err = range []error{initramFsErr, dracutErr, mkinitcpioErr} {
		if err == nil {
			continue
		}
		// we know this error isn't ErrNotExist, so we should throw it and exit
		log.Fatalf("Failed to stat file: %s", err)
	}

	switch {
	case initramFsExists:
		// Copy initramfs-tools module to system
		configs.CopyToSystem(config.IsRoot, initramfsFile, "/etc/initramfs-tools/modules")

		// Copy the modules file to /etc/modules
		configs.CopyToSystem(config.IsRoot, config.Path.ETCMODULES, "/etc/modules")

		if config.HasDuplicateDeviceIds {
			for configPath, sysPath := range config.EarlyBindFilePaths {
				configs.CopyToSystem(config.IsRoot, configPath, sysPath)
			}
		}

		if err = command.ExecAndLogSudo(config.IsRoot, true, "update-initramfs", "-u"); err != nil {
			log.Fatalf("Failed to update initramfs: %s", err)
		}

	case dracutExists:
		// Copy dracut config to /etc/dracut.conf.d/vfio
		configs.CopyToSystem(config.IsRoot, dracutFile, "/etc/dracut.conf.d/vfio")

		if config.HasDuplicateDeviceIds {
			moduleSysPath := strings.Replace(config.Path.DRACUTMODULE, "config", "", 1)
			if err := command.ExecAndLogSudo(config.IsRoot, false, "mkdir", "-p", moduleSysPath); err != nil {
				log.Fatalf("Failed to create dracut module directory: %s", err)
			}

			for configPath, sysPath := range config.EarlyBindFilePaths {
				configs.CopyToSystem(config.IsRoot, configPath, sysPath)
			}
		}

		// Get systeminfo
		sysinfo := uname.New()

		if err = command.ExecAndLogSudo(config.IsRoot, true, "dracut", "-f", "-v", "--kver", sysinfo.Release); err != nil {
			log.Fatalf("Failed to update initramfs: %s", err)
		}

	case mkinitcpioExists:
		// Copy dracut config to /etc/dracut.conf.d/vfio
		configs.CopyToSystem(config.IsRoot, config.Path.MKINITCPIO, "/etc/mkinitcpio.conf")

		if config.HasDuplicateDeviceIds {
			for configPath, sysPath := range config.EarlyBindFilePaths {
				configs.CopyToSystem(config.IsRoot, configPath, sysPath)
			}
		}

		if err = command.ExecAndLogSudo(config.IsRoot, true, "mkinitcpio", "-P"); err != nil {
			log.Fatalf("Failed to update initramfs: %s", err)
		}

	}

	// Make sure prompt end up on next line
	fmt.Print("\n")
}
