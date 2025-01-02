package configs

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

// Set_Dracut writes a dracut configuration file for `/etc/dracut.conf.d/`.
func Set_Dracut(config *Config) {
	// Set the dracut config file
	dracutConf := fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT)

	// If the file already exists then delete it
	if exists, _ := fileio.FileExist(dracutConf); exists {
		_ = os.Remove(dracutConf)
	}

	// Write to logger
	logger.Printf("Writing to %s:\nforce_drivers+=\" %s \"\n", dracutConf, strings.Join(vfio_modules(), " "))

	// Write the dracut config file
	fileio.AppendContent(fmt.Sprintf("force_drivers+=\" %s \"\n", strings.Join(vfio_modules(), " ")), dracutConf)

	// Get the current kernel arguments we have generated
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// If the kernel argument is not already in the file
	if !strings.Contains(kernel_args, "rd.driver.pre=vfio-pci") {
		// Add to our kernel arguments file that vfio_pci should load early (dracut does this using kernel arguments)
		fileio.AppendContent(" rd.driver.pre=vfio-pci", config.Path.CMDLINE)
	}

	// Make a backup of dracutConf if there is one there
	backupFile(strings.Replace(dracutConf, "config", "", 1))

	if config.HasDuplicateDeviceIds {
		setDracutEarlyBinds(config)
	}
}

func setDracutEarlyBinds(config *Config) {
	err := os.MkdirAll(config.Path.DRACUTMODULE, os.ModePerm)
	common.ErrorCheck(err, "Error, could not create dracut module config directory")
	confToSystemPathRe := regexp.MustCompile(`^config`)

	earlyBindScriptConfigPath := path.Join(config.Path.DRACUTMODULE, "early-vfio-bind.sh")
	earlyBindScriptSysPath := confToSystemPathRe.ReplaceAllString(earlyBindScriptConfigPath, "")
	config.EarlyBindFilePaths[earlyBindScriptConfigPath] = earlyBindScriptSysPath
	if exists, _ := fileio.FileExist(earlyBindScriptConfigPath); exists {
		_ = os.Remove(earlyBindScriptConfigPath)
	}

	logger.Printf("Writing to early bind script to %s", earlyBindScriptConfigPath)
	vfioBindScript := fmt.Sprintf(`#!/bin/bash
DEVS="%s"

for DEV in $DEVS; do
	echo "vfio-pci" > /sys/bus/pci/devices/$DEV/driver_override
done

# Load the vfio-pci module
modprobe -i vfio-pci`, strings.Join(config.Gpu_Addresses, " "))

	fileio.AppendContent(vfioBindScript, earlyBindScriptConfigPath)
	err = os.Chmod(earlyBindScriptConfigPath, 0755)
	common.ErrorCheck(err, fmt.Sprintf("Error, could not chmod %s", earlyBindScriptConfigPath))

	dracutModuleConfigPath := path.Join(config.Path.DRACUTMODULE, "module-setup.sh")
	dracutModuleSysPath := confToSystemPathRe.ReplaceAllString(dracutModuleConfigPath, "")
	config.EarlyBindFilePaths[dracutModuleConfigPath] = dracutModuleSysPath
	if exists, _ := fileio.FileExist(dracutModuleConfigPath); exists {
		_ = os.Remove(dracutModuleConfigPath)
	}

	logger.Printf("Writing to dracut early bind config to %s", dracutModuleConfigPath)
	dracutConfig := fmt.Sprintf(`#!/bin/bash
check() {
	return 0
}

depends() {
	return 0
}

install() {
	inst_hook pre-trigger 90 "$moddir/%s"
}`, path.Base(earlyBindScriptSysPath))

	fileio.AppendContent(dracutConfig, dracutModuleConfigPath)
	err = os.Chmod(dracutModuleConfigPath, 0755)
	common.ErrorCheck(err, fmt.Sprintf("Error, could not chmod %s", dracutModuleConfigPath))
}
