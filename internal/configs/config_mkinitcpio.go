package configs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

func Set_Mkinitcpio() {
	// Get the config struct
	config := GetConfig()

	// Make sure we start from scratch by deleting any old file
	if fileio.FileExist(config.Path.MKINITCPIO) {
		os.Remove(config.Path.MKINITCPIO)
	}

	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)
	sysfile := syspath_re.ReplaceAllString(config.Path.MKINITCPIO, "")

	// Make a regex to find the modules line
	module_line_re := regexp.MustCompile(`^MODULES=`)
	modules_re := regexp.MustCompile(`MODULES=\((.+?)\)`)
	vfio_modules_re := regexp.MustCompile(`(vfio_iommu_type1|vfio_pci|vfio_virqfd|vfio|vendor-reset)`)

	// Read the mkinitcpio file
	mkinitcpio_content := fileio.ReadLines(sysfile)

	for _, line := range mkinitcpio_content {
		// If we are at the line starting with MODULES=
		if module_line_re.MatchString(line) {
			// Get the current modules
			currentmodules := strings.Split(modules_re.ReplaceAllString(line, "${1}"), " ")

			// Get the vfio modules we need to use
			modules := vfio_modules()

			// If vendor-reset is in the current modules
			if strings.Contains(line, "vendor-reset") {
				// Add vendor-reset first
				modules = append([]string{"vendor-reset"}, modules...)
			}

			// Loop through current modules and add anything that isnt vfio or vendor-reset related
			for _, v := range currentmodules {
				// If what we find is not a vfio module or vendor-reset module
				if !vfio_modules_re.MatchString(v) {
					// Add module to module list
					modules = append(modules, v)
				}
			}

			// Write the modules line we generated
			fileio.AppendContent(fmt.Sprintf("MODULES=(%s)\n", strings.Join(modules, " ")), config.Path.MKINITCPIO)
		} else {
			// Else just write the line to the config
			fileio.AppendContent(fmt.Sprintf("%s\n", line), config.Path.MKINITCPIO)
		}
	}
}
