package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

// This function adds the disable vfio video output on host option to the config
// The function will use the given int as the value for the option
func DisableVFIOVideo(i int) {
	// Get the config
	config := GetConfig()

	// Write to logger
	logger.Printf("Adding vfio_pci.disable_vga=%v to %s", i, config.Path.CMDLINE)

	// Get the current kernel arguments we have generated
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// If the kernel argument is already in the file
	if strings.Contains(kernel_args, "vfio_pci.disable_vga") {
		// Remove the old file
		err := os.Remove(config.Path.CMDLINE)
		errorcheck.ErrorCheck(err, fmt.Sprintf("Could not rewrite %s", config.Path.CMDLINE))

		// Enable or disable the VGA based on our given value
		if i == 0 {
			kernel_args = strings.Replace(kernel_args, "vfio_pci.disable_vga=1", "vfio_pci.disable_vga=0", 1)

		} else {
			kernel_args = strings.Replace(kernel_args, "vfio_pci.disable_vga=0", "vfio_pci.disable_vga=1", 1)
		}

		// Rewrite the kernel_args file
		fileio.AppendContent(kernel_args, config.Path.CMDLINE)
	} else {
		// Add to the kernel arguments that we want to disable VFIO video output on the host
		fileio.AppendContent(
			fmt.Sprintf(
				" vfio_pci.disable_vga=%v", i,
			),
			config.Path.CMDLINE,
		)
	}
}
