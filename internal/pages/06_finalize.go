package pages

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	isRoot := os.Getuid() == 0

	config.IsRoot = isRoot

	finalizeNotice(isRoot)

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
		log.Fatalf(err.Error())
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
		output = configs.Set_KernelStub()
		fmt.Printf("%s\n", output)

	} else if config.Bootloader == "grubby" {
		// Write to logger
		logger.Printf("Configuring bootloader using grubby\n")

		// Configure kernelstub
		output = configs.Set_Grubby()
		fmt.Printf("%s\n", output)

	} else if config.Bootloader == "grub2" {
		// Write to logger
		logger.Printf("Applying grub2 changes\n")
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

	execAndLogSudo := func(cmd string) {
		if !config.IsRoot && !strings.HasPrefix(cmd, "sudo") {
			cmd = fmt.Sprintf("sudo %s", cmd)
		}
		// Write to logger
		logger.Printf("Executing: %s\n", cmd)

		// Update initramfs
		fmt.Printf("Executing: %s\nSee debug.log for detailed output\n", cmd)
		cs := strings.Fields(cmd)
		r := exec.Command(cs[0], cs[1:]...)

		cmd_out, _ := r.CombinedOutput()

		// Write to logger
		logger.Printf(string(cmd_out) + "\n")
	}

	// Copy the config files for the system we have
	initramfsFile := fmt.Sprintf("%s/modules", config.Path.INITRAMFS)
	dracutFile := fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT)
	switch {
	case fileio.FileExist(initramfsFile):
		// Copy initramfs-tools module to system
		output = configs.CopyToSystem(initramfsFile, "/etc/initramfs-tools/modules")
		fmt.Printf("%s\n", output)

		// Copy the modules file to /etc/modules
		output = configs.CopyToSystem(config.Path.ETCMODULES, "/etc/modules")
		fmt.Printf("%s\n", output)

		execAndLogSudo("update-initramfs -u")

	case fileio.FileExist(dracutFile):
		// Copy dracut config to /etc/dracut.conf.d/vfio
		output = configs.CopyToSystem(dracutFile, "/etc/dracut.conf.d/vfio")
		fmt.Printf("%s\n", output)

		// Get systeminfo
		sysinfo := uname.New()

		execAndLogSudo(fmt.Sprintf("dracut -f -v --kver %s\n", sysinfo.Release))

	case fileio.FileExist(config.Path.MKINITCPIO):
		// Copy dracut config to /etc/dracut.conf.d/vfio
		output = configs.CopyToSystem(config.Path.MKINITCPIO, "/etc/mkinitcpio.conf")
		fmt.Printf("%s\n", output)

		execAndLogSudo("mkinitcpio -P")

	}

	// Make sure prompt end up on next line
	fmt.Print("\n")
}
