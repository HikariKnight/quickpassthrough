package configs

import (
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/klauspost/cpuid/v2"
)

// This function just adds what bootloader the system has to our config.bootloader value
// Preference is given to kernelstub because it is WAY easier to safely edit compared to grub2
func getBootloader(config *Config) {
	// Check what bootloader handler we are using
	// Check for grub-mkconfig
	_, err := command.Run("which", "grub-mkconfig")
	if err == nil {
		// Mark bootloader as grub2
		config.Bootloader = "grub2"
	}

	// Check for grubby (used by fedora)
	_, err = command.Run("which", "grubby")
	if err == nil {
		// Mark it as unknown as i do not support it yet
		config.Bootloader = "unknown"
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
func set_Cmdline() {
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

	// If the config folder for dracut exists in our configs
	if fileio.FileExist(config.Path.DRACUT) {
		// Add an extra kernel argument needed for dracut users
		fileio.AppendContent(" rd.driver.pre=vfio_pci", config.Path.CMDLINE)
	}
}
