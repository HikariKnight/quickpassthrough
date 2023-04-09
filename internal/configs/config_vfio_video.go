package configs

import (
	"fmt"

	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

func DisableVFIOVideo(i int) {
	// Get the config
	config := GetConfig()

	// Add to the kernel arguments that we want to disable VFIO video output on the host
	fileio.AppendContent(
		fmt.Sprintf(
			" vfio_pci.disable_vga=%v", i,
		),
		config.Path.CMDLINE,
	)
}
