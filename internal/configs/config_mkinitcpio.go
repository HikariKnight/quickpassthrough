package configs

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

// Set_Mkinitcpio copies the content of /etc/mkinitcpio.conf to the config folder and does an inline replace/insert on the MODULES=() line
func Set_Mkinitcpio(config *Config) {
	// Make sure we start from scratch by deleting any old file
	if exists, _ := fileio.FileExist(config.Path.MKINITCPIO); exists {
		_ = os.Remove(config.Path.MKINITCPIO)
	}

	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)
	sysfile := syspath_re.ReplaceAllString(config.Path.MKINITCPIO, "")

	// Make a regex to find the modules line
	module_line_re := regexp.MustCompile(`^MODULES=`)
	hooks_line_re := regexp.MustCompile(`^HOOKS=`)
	modules_re := regexp.MustCompile(`MODULES=\((.*)\)`)
	vfio_modules_re := regexp.MustCompile(`(vfio_iommu_type1|vfio_pci|vfio_virqfd|vfio|vendor-reset)`)

	// Read the mkinitcpio file
	mkinitcpio_content := fileio.ReadLines(sysfile)

	// Write to logger
	logger.Printf("Read %s:\n%s\n", sysfile, strings.Join(mkinitcpio_content, "\n"))

	for _, line := range mkinitcpio_content {
		// If we are at the line starting with MODULES=
		if module_line_re.MatchString(line) {
			// Get the current modules
			currentmodules := strings.Split(modules_re.ReplaceAllString(line, "${1}"), " ")

			// Get the vfio modules we need to use
			modules := vfio_modules()

			// If vendor-reset is in the current modules
			if strings.Contains(line, "vendor-reset") {
				// Write to logger
				logger.Printf("vendor-reset module detected in %s\nMaking sure it will be loaded before vfio\n", sysfile)

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

			// Write to logger
			logger.Printf("Replacing line in %s:\n%s\nWith:\nMODULES=(%s)\n", config.Path.MKINITCPIO, line, strings.Join(modules, " "))

			// Write the modules line we generated
			fileio.AppendContent(fmt.Sprintf("MODULES=(%s)\n", strings.Join(modules, " ")), config.Path.MKINITCPIO)
		} else if config.HasDuplicateDeviceIds && hooks_line_re.MatchString(line) {
			setMkinitcpioEarlyBinds(config, line)
		} else {
			// Else just write the line to the config
			fileio.AppendContent(fmt.Sprintf("%s\n", line), config.Path.MKINITCPIO)
		}
	}
}

func setMkinitcpioEarlyBinds(config *Config, hooksLine string) {
	err := os.MkdirAll(config.Path.MKINITCPIOHOOKS, os.ModePerm)
	common.ErrorCheck(err, "Error, could not create mkinitcpio hook config directory")
	confToSystemPathRe := regexp.MustCompile(`^config`)

	earlyBindHookConfigPath := path.Join(config.Path.MKINITCPIOHOOKS, "early-vfio-bind")
	earlyBindHookSysPath := confToSystemPathRe.ReplaceAllString(earlyBindHookConfigPath, "")
	config.EarlyBindFilePaths[earlyBindHookConfigPath] = earlyBindHookSysPath
	if exists, _ := fileio.FileExist(earlyBindHookConfigPath); exists {
		_ = os.Remove(earlyBindHookConfigPath)
	}

	logger.Printf("Writing to early bind hook to %s", earlyBindHookConfigPath)
	vfioBindHook := fmt.Sprintf(`#!/bin/bash
run_hook() {
	DEVS="%s"

	for DEV in $DEVS; do
		echo "vfio-pci" > /sys/bus/pci/devices/$DEV/driver_override
	done

	# Load the vfio-pci module
	modprobe -i vfio-pci
}`, strings.Join(config.Gpu_Addresses, " "))

	fileio.AppendContent(vfioBindHook, earlyBindHookConfigPath)
	err = os.Chmod(earlyBindHookConfigPath, 0755)
	common.ErrorCheck(err, fmt.Sprintf("Error, could not chmod %s", earlyBindHookConfigPath))

	hooksString := strings.Trim(strings.Split(hooksLine, "=")[1], "()")
	hooks := strings.Split(hooksString, " ")
	customHook := "early-vfio-bind"
	if !slices.Contains(hooks, customHook) {
		hooks = append(hooks, customHook)
	}

	// Write to logger
	logger.Printf("Replacing line in %s:\n%s\nWith:\nHOOKS=(%s)\n", config.Path.MKINITCPIO, hooksLine, strings.Join(hooks, " "))

	// Write the modules line we generated
	fileio.AppendContent(fmt.Sprintf("HOOKS=(%s)\n", strings.Join(hooks, " ")), config.Path.MKINITCPIO)
}
