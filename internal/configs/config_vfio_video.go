package configs

import "github.com/HikariKnight/quickpassthrough/pkg/fileio"

func DisableVFIOVideo() {
	// Get the config
	config := GetConfig()

	// Add to the kernel arguments that we want to disable VFIO video output on the host
	fileio.AppendContent(" vfio_pci.disable_vga=1", config.Path.CMDLINE)
}
