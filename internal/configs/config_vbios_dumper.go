package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
)

// Generates a script file named dump_vbios.sh and places it inside the utils folder.
// This script has to be run without a display manager or display server running
func GenerateVBIOSDumper(vbios_path string) {
	// Get the config directories
	config := GetConfig()

	// Get the program directory
	exe, _ := os.Executable()
	scriptdir := filepath.Dir(exe)

	// If we are using go run use the working directory instead
	if strings.Contains(scriptdir, "/tmp/go-build") {
		scriptdir, _ = os.Getwd()
	}

	vbios_script_template := fmt.Sprint(
		"#!/bin/bash\n",
		"# THIS FILE IS AUTO GENERATED!\n",
		"# IF YOU HAVE CHANGED GPU, PLEASE RE-RUN QUICKPASSTHROUGH!\n",
		"mkdir -p \"%s\"\n",
		"echo Attempting to enable reading from rom\n",
		"echo 1 | sudo tee %s\n",
		"echo\n",
		"echo Attempting to dump VBIOS\n",
		"sudo bash -c \"cat %s\" > %s/%s/vfio_card.rom || echo \"\nFailed to dump the VBIOS, in most cases a reboot can fix this.\nOr you have to bind the gpu to the vfio-pci driver, reboot the machine and try dumping again.\nIf that still fails, you might find your VBIOS at: https://www.techpowerup.com/vgabios/\n\"\n",
		"file \"%s/%s/vfio_card.rom\"\n",
		"echo\n",
		"echo Attempting to disable reading from rom \\(cleanup\\)\n",
		"echo 0 | sudo tee %s\n",
	)

	vbios_script := fmt.Sprintf(
		vbios_script_template,
		config.Path.QUICKEMU,
		vbios_path,
		vbios_path,
		scriptdir,
		config.Path.QUICKEMU,
		scriptdir,
		config.Path.QUICKEMU,
		vbios_path,
	)

	// Make the script file
	scriptfile, err := os.Create("utils/dump_vbios.sh")
	errorcheck.ErrorCheck(err, "Cannot create file \"utils/dump_vbios.sh\"")
	defer scriptfile.Close()

	// Make the script executable
	scriptfile.Chmod(0775)
	errorcheck.ErrorCheck(err, "Could not change permissions of \"utils/dump_vbios.sh\"")

	// Write to logger
	logger.Printf("Writing utils/dump_vbios.sh\n")

	// Write the script
	scriptfile.WriteString(vbios_script)
}
