package configs

import (
	"errors"
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
	bootloader string
	cpuvendor  string
	path       *Path
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
	config := &Config{}
	config.path = GetConfigPaths()

	// Set default value for bootloader
	config.bootloader = "unknown"

	// Detect the bootloader we are using
	getBootloader(config)

	// Detect the cpu vendor
	config.cpuvendor = cpuid.CPU.VendorString

	return config
}

func InitConfigs() {
	config := GetConfig()

	// Add all directories we need into a stringlist
	dirs := []string{
		config.path.MODPROBE,
		config.path.INITRAMFS,
		config.path.DEFAULT,
		config.path.DRACUT,
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
		if _, err := os.Stat(syspath); !errors.Is(err, os.ErrNotExist) {
			// Create the directories for our configs
			err := os.MkdirAll(confpath, os.ModePerm)
			errorcheck.ErrorCheck(err)
		}
	}

	// Add all files we need to a stringlist
	files := []string{
		config.path.ETCMODULES,
		config.path.MKINITCPIO,
		fmt.Sprintf("%s/modules", config.path.INITRAMFS),
		fmt.Sprintf("%s/grub", config.path.DEFAULT),
	}

	for _, conffile := range files {
		// Get the system file path
		sysfile := syspath_re.ReplaceAllString(conffile, "")

		// If the file exists
		if _, err := os.Stat(sysfile); !errors.Is(err, os.ErrNotExist) {
			// Create the directories for our configs
			file, err := os.Create(conffile)
			errorcheck.ErrorCheck(err)
			// Close the file so we can edit it
			file.Close()
		}

		// If we now have a config that exists
		if _, err := os.Stat(conffile); !errors.Is(err, os.ErrNotExist) {
			switch conffile {
			case config.path.ETCMODULES:
				// Read the header
				header := initramfs_readHeader(4, sysfile)
				fileio.AppendContent(header, conffile)

				// Add the modules to the config file
				initramfs_addModules(conffile)
			case config.path.INITRAMFS:
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
