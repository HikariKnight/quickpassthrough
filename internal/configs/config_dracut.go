package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

// This function writes a dracut configuration file for /etc/dracut.conf.d/
func Set_Dracut() {
	config := GetConfig()

	// Set the dracut config file
	dracutConf := fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT)

	// If the file already exists then delete it
	if fileio.FileExist(dracutConf) {
		os.Remove(dracutConf)
	}

	// Write to logger
	logger.Printf("Writing to %s:\nadd_drivers+=\" %s \"", dracutConf, strings.Join(vfio_modules(), " "))

	// Write the dracut config file
	fileio.AppendContent(fmt.Sprintf("add_drivers+=\" %s \"\n", strings.Join(vfio_modules(), " ")), dracutConf)

	// Get the current kernel arguments we have generated
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// If the kernel argument is not already in the file
	if !strings.Contains(kernel_args, "rd.driver.pre=vfio_pci") {
		// Add to our kernel arguments file that vfio_pci should load early (dracut does this using kernel arguments)
		fileio.AppendContent(" rd.driver.pre=vfio_pci", config.Path.CMDLINE)
	}

	// Make a backup of dracutConf if there is one there
	backupFile(strings.Replace(dracutConf, "config", "", 1))
}
