package configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/klauspost/cpuid/v2"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/HikariKnight/quickpassthrough/pkg/uname"
)

type Path struct {
	CMDLINE      string
	MODPROBE     string
	INITRAMFS    string
	ETCMODULES   string
	DEFAULT      string
	QEMU         string
	DRACUT       string
	DRACUTMODULE string
	MKINITCPIO   string
}

type Config struct {
	Bootloader            string
	Cpuvendor             string
	Path                  *Path
	Gpu_Group             string
	Gpu_IDs               []string
	Gpu_Addresses         []string
	EarlyBindFilePaths    map[string]string
	IsRoot                bool
	HasDuplicateDeviceIds bool
}

// GetConfigPaths retrieves the path to all the config files.
func GetConfigPaths() *Path {
	Paths := &Path{
		CMDLINE:      "config/kernel_args",
		MODPROBE:     "config/etc/modprobe.d",
		INITRAMFS:    "config/etc/initramfs-tools",
		ETCMODULES:   "config/etc/modules",
		DEFAULT:      "config/etc/default",
		QEMU:         "config/qemu",
		DRACUT:       "config/etc/dracut.conf.d",
		DRACUTMODULE: "config/usr/lib/dracut/modules.d/90early-vfio-bind",
		MKINITCPIO:   "config/etc/mkinitcpio.conf",
	}

	return Paths
}

// GetConfig retrieves all the configs and returns the struct.
func GetConfig() *Config {
	config := &Config{
		Bootloader:            "unknown",
		Cpuvendor:             cpuid.CPU.VendorString,
		Path:                  GetConfigPaths(),
		Gpu_Group:             "",
		Gpu_IDs:               []string{},
		Gpu_Addresses:         []string{},
		EarlyBindFilePaths:    map[string]string{},
		HasDuplicateDeviceIds: false,
	}

	// Detect the bootloader we are using
	getBootloader(config)

	return config
}

// InitConfigs constructs the empty config files and folders based on what exists on the system
func InitConfigs() {
	config := GetConfig()

	// Add all directories we need into a stringlist
	dirs := []string{
		config.Path.MODPROBE,
		config.Path.INITRAMFS,
		config.Path.DEFAULT,
		config.Path.DRACUT,
		config.Path.DRACUTMODULE,
	}

	// Remove old config
	if err := os.RemoveAll("config"); err != nil && !errors.Is(err, os.ErrNotExist) {

		// won't be called if the error is ErrNotExist
		common.ErrorCheck(err, "\nError removing old config")
	}

	// Make the config folder
	if err := os.Mkdir("config", os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		// won't be called if the error is ErrExist
		common.ErrorCheck(err, "\nError making config folder")
	}

	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)

	// For each directory
	for _, confpath := range dirs {
		// Get the system path
		syspath := syspath_re.ReplaceAllString(confpath, "")

		exists, err := fileio.FileExist(syspath)

		// If we received an error that is not ErrNotExist
		if err != nil {
			common.ErrorCheck(err, "\nError checking for directory: "+syspath)
			continue // note: unreachable due to ErrorCheck calling fatal
		}

		// If the path exists
		if exists {
			// Write to log
			logger.Printf(
				"%s found on the system\n"+
					"Creating %s\n",
				syspath,
				confpath,
			)

			// Make a backup directory
			makeBackupDir(syspath)

			// Create the directories for our configs
			if err = os.MkdirAll(confpath, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
				common.ErrorCheck(err, "\nError making directory: "+confpath)
				return // note: unreachable due to ErrorCheck calling fatal
			}
		}
	}

	// Add all files we need to a stringlist
	files := []string{
		config.Path.ETCMODULES,
		config.Path.MKINITCPIO,
		fmt.Sprintf("%s/modules", config.Path.INITRAMFS),
		fmt.Sprintf("%s/grub", config.Path.DEFAULT),
	}

	// If we are using grubby
	if config.Bootloader == "grubby" {
		// Do not create an empty /etc/default/grub file
		files = files[:len(files)-1]
	}

	for _, conffile := range files {
		// Get the system file path
		sysfile := syspath_re.ReplaceAllString(conffile, "")

		// If the file exists
		exists, err := fileio.FileExist(sysfile)

		// If we received an error that is not ErrNotExist
		if err != nil {
			common.ErrorCheck(err, "\nError checking for file: "+sysfile)
			continue // note: unreachable due to ErrorCheck calling fatal
		}

		if exists {
			// Write to log
			logger.Printf(
				"%s found on the system\n"+
					"Creating %s\n",
				sysfile,
				conffile,
			)

			// Create the directories for our configs
			file, err := os.Create(conffile)
			common.ErrorCheck(err)
			// Close the file so we can edit it
			_ = file.Close()

			// Backup the sysfile if we do not have a backup
			backupFile(sysfile)
		}

		exists, err = fileio.FileExist(conffile)
		if err != nil {
			common.ErrorCheck(err, "\nError checking for file: "+conffile)
			continue // note: unreachable
		}

		// If we now have a config that exists
		if exists {
			switch conffile {
			case config.Path.ETCMODULES:
				// Write to logger
				logger.Printf("Getting the header (if it is there) from %s\n", conffile)

				// Read the header
				header := initramfs_readHeader(4, sysfile)
				fileio.AppendContent(header, conffile)

				// Add the modules to the config file
				initramfs_addModules(conffile)
			case fmt.Sprintf("%s/modules", config.Path.INITRAMFS):
				// Write to logger
				logger.Printf("Getting the header (if it is there) from %s\n", conffile)

				// Read the header
				header := initramfs_readHeader(11, sysfile)
				fileio.AppendContent(header, conffile)

				// Add the modules to the config file
				initramfs_addModules(conffile)
			}
		}
	}
}

// Returns a list of modules used for vfio based on the systems kernel version
func vfio_modules() []string {
	// Make the list of modules
	modules := []string{
		"vfio_pci",
		"vfio",
		"vfio_iommu_type1",
	}

	// If we are on a kernel older than 6.2
	sysinfo := uname.New()
	kernel_re := regexp.MustCompile(`^(6\.1|6\.0|[1-5]\.\d{1,2})\.`)
	if kernel_re.MatchString(sysinfo.Kernel) {
		// Write to the debug log
		logger.Printf("Linux kernel version %s detected!\nIncluding vfio_virqfd module\n", sysinfo.Kernel)

		// Include the vfio_virqfd module
		// NOTE: this driver was merged into the vfio module in 6.2
		modules = append(modules, "vfio_virqfd")
	}

	// Return the modules
	return modules
}

func backupFile(source string) {
	// Make a destination path
	dest := fmt.Sprintf("backup%s", source)

	configExists, configFileError := fileio.FileExist(fmt.Sprintf("config%s", source))
	sysExists, sysFileError := fileio.FileExist(source)
	destExists, destFileError := fileio.FileExist(dest)

	// If we received an error that is not ErrNotExist on any of the files
	for _, err := range []error{configFileError, sysFileError, destFileError} {
		if err != nil {
			common.ErrorCheck(configFileError, "\nError checking for file: "+source)
			return // note: unreachable
		}
	}

	switch {
	// If the file exists in the config but not on the system it is a file we make
	case configExists && !sysExists:
		// Create the blank file so that a copy of the backup folder to /etc
		file, err := os.Create(dest)
		common.ErrorCheck(err, "Error creating file %s\n", dest)
		_ = file.Close()

		// If a backup of the file does not exist
	case sysExists && !destExists:
		// Write to the logger
		logger.Printf("No first time backup of %s detected.\nCreating a backup at %s\n", source, dest)

		// Copy the file
		fileio.FileCopy(source, dest)
	}

}

func makeBackupDir(dest string) {
	// If a backup directory does not exist
	exists, err := fileio.FileExist("backup/")
	if err != nil {
		// If we received an error that is not ErrNotExist
		common.ErrorCheck(err, "Error checking for backup/ folder")
		return // note: unreachable
	}

	if !exists {
		// Write to the logger
		logger.Printf("Backup directory does not exist!\nCreating backup directory for first run backup")
	}

	// Make the empty directories
	if err = os.MkdirAll(fmt.Sprintf("backup/%s", dest), os.ModePerm); errors.Is(err, os.ErrExist) {
		// ignore if the directory already exists
		err = nil
	}
	// will return without incident if there's no error
	common.ErrorCheck(err, "Error making backup/ folder")
}

// CopyToSystem copies a file to the system.
func CopyToSystem(isRoot bool, conffile, sysfile string) {
	// Since we should be elevated with our sudo token we will copy with cp
	// (using built in functions will not work as we are running as the normal user)

	// ExecAndLogSudo will write to the logger, so just print here
	fmt.Printf("Copying: %s to %s\n", conffile, sysfile)

	if isRoot {
		logger.Printf("Copying %s to %s\n", conffile, sysfile)
		fmt.Printf("Copying %s to %s\n", conffile, sysfile)
		fDat, err := os.ReadFile(conffile)
		common.ErrorCheck(err, fmt.Sprintf("Failed to read %s", conffile))
		err = os.WriteFile(sysfile, fDat, 0644)
		common.ErrorCheck(err, fmt.Sprintf("Failed to write %s", sysfile))
		logger.Printf("Copied %s to %s\n", conffile, sysfile)
		return
	}

	if !filepath.IsAbs(conffile) {
		conffile, _ = filepath.Abs(conffile)
	}

	err := command.ExecAndLogSudo(isRoot, false, "cp", "-v", conffile, sysfile)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// [command.ExecAndLogSudo] will log the command's output
	common.ErrorCheck(err, fmt.Sprintf("Failed to copy %s to %s:\n%s", conffile, sysfile, errMsg))

	// ---------------------------------------------------------------------------------
	// note that if we failed the error check, the following will not appear in the log!
	// this is because the [common.ErrorCheck] function will call [log.Fatalf] and exit
	// ---------------------------------------------------------------------------------

	logger.Printf("Copied %s to %s\n", conffile, sysfile)
}
