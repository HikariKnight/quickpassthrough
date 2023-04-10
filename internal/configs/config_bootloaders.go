package configs

import (
	"fmt"
	"strings"

	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/klauspost/cpuid/v2"
)

// This function just adds what bootloader the system has to our config.bootloader value
// Preference is given to kernelstub because it is WAY easier to safely edit compared to grub2
func getBootloader(config *Config) {
	// Check what bootloader handler we are using
	// Check for grub-mkconfig
	_, err := command.Run("which", "grub2-mkconfig")
	if err == nil {
		// Mark bootloader as grub2
		config.Bootloader = "grub2"
	}

	// Check for grubby (used by fedora)
	_, err = command.Run("which", "grubby")
	if err == nil {
		// Mark it as unknown as i do not support it yet
		config.Bootloader = "grubby"
	}

	// Check for kernelstub (used by pop os)
	_, err = command.Run("which", "kernelstub")
	if err == nil {
		config.Bootloader = "kernelstub"
	}
}

// This function adds the default kernel arguments we want to the config/cmdline file
// This gives us a file we can read all the kernel arguments this system needs
// in case of an unknown bootloader
func Set_Cmdline(gpu_IDs []string) {
	// Get the system info
	cpuinfo := cpuid.CPU

	// Get the configs
	config := GetConfig()

	// Write the file containing our kernel arguments to feed the bootloader
	fileio.AppendContent("iommu=pt", config.Path.CMDLINE)

	// Write the argument based on which cpu the user got
	switch cpuinfo.VendorString {
	case "AuthenticAMD":
		fileio.AppendContent(" amd_iommu=on", config.Path.CMDLINE)
	case "GenuineIntel":
		fileio.AppendContent(" intel_iommu=on", config.Path.CMDLINE)
	}

	// Add the GPU ids for vfio to the kernel arguments
	fileio.AppendContent(fmt.Sprintf(" vfio_pci.ids=%s", strings.Join(gpu_IDs, ",")), config.Path.CMDLINE)
}

// TODO: write functions to configure grub
// TODO3: if unknown bootloader, tell user what to add a kernel arguments

// Configures systemd-boot using kernelstub
func Set_KernelStub() {
	// Get the config
	config := GetConfig()

	// Get the kernel args
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// Write to logger
	logger.Printf("Running command:\nsudo kernelstub -a \"%s\"", kernel_args)

	// Run the command
	command.Run("sudo", "kernelstub", "-a", kernel_args)
}

// Configures grub2 or systemd-boot using grubby
func Set_Grubby() {
	// Get the config
	config := GetConfig()

	// Get the kernel args
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// Write to logger
	logger.Printf("Running command:\nsudo grubby --update-kernel=ALL --args=\"%s\"", kernel_args)

	// Run the command
	command.Run("sudo", "grubby", "--update-kernel=ALL", fmt.Sprintf("--args=%s", kernel_args))
}
