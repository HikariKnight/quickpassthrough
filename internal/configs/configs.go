package configs

import (
	"fmt"
	"os"
	"regexp"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/HikariKnight/quickpassthrough/pkg/uname"
	"github.com/klauspost/cpuid/v2"
)

type Path struct {
	CMDLINE    string
	MODPROBE   string
	INITRAMFS  string
	ETCMODULES string
	DEFAULT    string
	QEMU       string
	DRACUT     string
	MKINITCPIO string
}

type Config struct {
	Bootloader string
	Cpuvendor  string
	Path       *Path
	Gpu_Group  string
	Gpu_IDs    []string
}

// Gets the path to all the config files
func GetConfigPaths() *Path {
	Paths := &Path{
		CMDLINE:    "config/kernel_args",
		MODPROBE:   "config/etc/modprobe.d",
		INITRAMFS:  "config/etc/initramfs-tools",
		ETCMODULES: "config/etc/modules",
		DEFAULT:    "config/etc/default",
		QEMU:       "config/qemu",
		DRACUT:     "config/etc/dracut.conf.d",
		MKINITCPIO: "config/etc/mkinitcpio.conf",
	}

	return Paths
}

// Gets all the configs and returns the struct
func GetConfig() *Config {
	config := &Config{
		Bootloader: "unknown",
		Cpuvendor:  cpuid.CPU.VendorString,
		Path:       GetConfigPaths(),
		Gpu_Group:  "",
		Gpu_IDs:    []string{},
	}

	// Detect the bootloader we are using
	getBootloader(config)

	return config
}

// Constructs the empty config files and folders based on what exists on the system
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

	// If we are using grubby
	if config.Bootloader == "grubby" {
		// Do not create an empty /etc/default/grub file
		files = files[:len(files)-1]
	}

	for _, conffile := range files {
		// Get the system file path
		sysfile := syspath_re.ReplaceAllString(conffile, "")

		// If the file exists
		if fileio.FileExist(sysfile) {
			// Write to log
			logger.Printf(
				"%s found on the system\n"+
					"Creating %s\n",
				sysfile,
				conffile,
			)

			// Create the directories for our configs
			file, err := os.Create(conffile)
			errorcheck.ErrorCheck(err)
			// Close the file so we can edit it
			file.Close()

			// Backup the sysfile if we do not have a backup
			backupFile(sysfile)
		}

		// If we now have a config that exists
		if fileio.FileExist(conffile) {
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
	kernel_re := regexp.MustCompile(`^(6\.1|6\.0|[1-5]\.)`)
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

	// If the file exists in the config but not on the system it is a file we make
	if fileio.FileExist(fmt.Sprintf("config%s", source)) && !fileio.FileExist(source) {
		// Create the blank file so that a copy of the backup folder to /etc
		file, err := os.Create(dest)
		errorcheck.ErrorCheck(err, "Error creating file %s\n", dest)
		file.Close()
	} else if !fileio.FileExist(dest) {
		// If a backup of the file does not exist
		// Write to the logger
		logger.Printf("No first time backup of %s detected.\nCreating a backup at %s\n", source, dest)

		// Copy the file
		fileio.FileCopy(source, dest)
	}

}

func makeBackupDir(dest string) {
	// If a backup directory does not exist
	if !fileio.FileExist("backup/") {
		// Write to the logger
		logger.Printf("Backup directory does not exist!\nCreating backup directory for first run backup")
	}

	// Make the empty directories
	err := os.MkdirAll(fmt.Sprintf("backup/%s", dest), os.ModePerm)
	errorcheck.ErrorCheck(err, "Error making backup/ folder")
}

// Copy a file to the system, make sure you have run command.Elevate() recently
func CopyToSystem(conffile, sysfile string) string {
	// Since we should be elevated with our sudo token we will copy with cp
	// (using built in functions will not work as we are running as the normal user)
	output, _ := command.Run("sudo", "cp", "-v", conffile, sysfile)

	// Clean the output
	clean_re := regexp.MustCompile(`\n`)
	clean_output := clean_re.ReplaceAllString(output[0], "")

	// Write output to logger
	logger.Printf("%s\n", clean_output)

	// Return the output
	return fmt.Sprintf("Copying: %s", clean_output)
}
