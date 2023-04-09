package configs

import (
	"fmt"
	"os"
	"regexp"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/klauspost/cpuid/v2"
)

type Path struct {
	CMDLINE    string
	MODPROBE   string
	INITRAMFS  string
	ETCMODULES string
	DEFAULT    string
	QUICKEMU   string
	DRACUT     string
	MKINITCPIO string
}

type Config struct {
	Bootloader string
	Cpuvendor  string
	Path       *Path
}

func GetConfigPaths() *Path {
	Paths := &Path{
		CMDLINE:    "config/cmdline",
		MODPROBE:   "config/etc/modprobe.d",
		INITRAMFS:  "config/etc/initramfs-tools",
		ETCMODULES: "config/etc/modules",
		DEFAULT:    "config/etc/default",
		QUICKEMU:   "config/quickemu",
		DRACUT:     "config/etc/dracut.conf.d",
		MKINITCPIO: "config/etc/mkinitcpio.conf",
	}

	return Paths
}

func GetConfig() *Config {
	config := &Config{
		Bootloader: "unknown",
		Cpuvendor:  cpuid.CPU.VendorString,
		Path:       GetConfigPaths(),
	}

	// Detect the bootloader we are using
	getBootloader(config)

	return config
}

func InitConfigs() {
	config := GetConfig()

	// Add all directories we need into a stringlist
	dirs := []string{
		config.Path.MODPROBE,
		config.Path.INITRAMFS,
		config.Path.DEFAULT,
		config.Path.DRACUT,
	}

	// Remove old config
	os.RemoveAll("config")

	// Make the config folder
	os.Mkdir("config", os.ModePerm)

	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)

	// For each directory
	for _, confpath := range dirs {
		// Get the system path
		syspath := syspath_re.ReplaceAllString(confpath, "")

		// If the path exists
		if fileio.FileExist(syspath) {
			// Create the directories for our configs
			err := os.MkdirAll(confpath, os.ModePerm)
			errorcheck.ErrorCheck(err)
		}
	}

	// Add all files we need to a stringlist
	files := []string{
		config.Path.ETCMODULES,
		config.Path.MKINITCPIO,
		fmt.Sprintf("%s/modules", config.Path.INITRAMFS),
		fmt.Sprintf("%s/grub", config.Path.DEFAULT),
	}

	for _, conffile := range files {
		// Get the system file path
		sysfile := syspath_re.ReplaceAllString(conffile, "")

		// If the file exists
		if fileio.FileExist(sysfile) {
			// Create the directories for our configs
			file, err := os.Create(conffile)
			errorcheck.ErrorCheck(err)
			// Close the file so we can edit it
			file.Close()
		}

		// If we now have a config that exists
		if fileio.FileExist(conffile) {
			switch conffile {
			case config.Path.ETCMODULES:
				// Read the header
				header := initramfs_readHeader(4, sysfile)
				fileio.AppendContent(header, conffile)

				// Add the modules to the config file
				initramfs_addModules(conffile)
			case fmt.Sprintf("%s/modules", config.Path.INITRAMFS):
				// Read the header
				header := initramfs_readHeader(11, sysfile)
				fileio.AppendContent(header, conffile)

				// Add the modules to the config file
				initramfs_addModules(conffile)
			}
		}
	}

	// Generate the kernel arguments
	set_Cmdline()
}
