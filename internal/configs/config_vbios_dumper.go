package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
)

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
		"echo 1 | sudo tee %s\n",
		"sudo bash -c \"cat %s\" > %s/%s/vfio_card.rom\n",
		"echo 0 | sudo tee %s\n",
	)

	vbios_script := fmt.Sprintf(
		vbios_script_template,
		vbios_path,
		vbios_path,
		scriptdir,
		config.Path.QUICKEMU,
		vbios_path,
	)

	scriptfile, err := os.Create("utils/dump_vbios.sh")
	errorcheck.ErrorCheck(err, "Cannot create file \"utils/dump_vbios.sh\"")
	defer scriptfile.Close()

	// Make the script executable
	scriptfile.Chmod(0775)
	errorcheck.ErrorCheck(err, "Could not change permissions of \"utils/dump_vbios.sh\"")

	// Write the script
	scriptfile.WriteString(vbios_script)
}
