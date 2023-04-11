package configs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/klauspost/cpuid/v2"
)

// This function just adds what bootloader the system has to our config.bootloader value
// Preference is given to kernelstub because it is WAY easier to safely edit compared to grub2
func getBootloader(config *Config) {
	// Check what bootloader handler we are using
	// Check for grub2-mkconfig
	_, err := command.Run("which", "grub2-mkconfig")
	if err == nil {
		// Mark bootloader as grub2
		config.Bootloader = "grub2"
	}

	// Check for grub2-mkconfig
	_, err = command.Run("which", "grub-mkconfig")
	if err == nil {
		// Mark bootloader as grub2
		config.Bootloader = "grub2"
	}

	// Check for grubby (used by fedora)
	_, err = command.Run("which", "grubby")
	if err == nil {
		// Mark it as unknown as i do not support it yet
		config.Bootloader = "grubby"
	}

	// Check for kernelstub (used by pop os)
	_, err = command.Run("which", "kernelstub")
	if err == nil {
		config.Bootloader = "kernelstub"
	}
}

// This function adds the default kernel arguments we want to the config/cmdline file
// This gives us a file we can read all the kernel arguments this system needs
// in case of an unknown bootloader
func Set_Cmdline(gpu_IDs []string) {
	// Get the system info
	cpuinfo := cpuid.CPU

	// Get the configs
	config := GetConfig()

	// Write the file containing our kernel arguments to feed the bootloader
	fileio.AppendContent("iommu=pt", config.Path.CMDLINE)

	// Write the argument based on which cpu the user got
	switch cpuinfo.VendorString {
	case "AuthenticAMD":
		fileio.AppendContent(" amd_iommu=on", config.Path.CMDLINE)
	case "GenuineIntel":
		fileio.AppendContent(" intel_iommu=on", config.Path.CMDLINE)
	}

	// Add the GPU ids for vfio to the kernel arguments
	fileio.AppendContent(fmt.Sprintf(" vfio_pci.ids=%s", strings.Join(gpu_IDs, ",")), config.Path.CMDLINE)
}

// Configures systemd-boot using kernelstub
func Set_KernelStub() string {
	// Get the config
	config := GetConfig()

	// Get the kernel args
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// Write to logger
	logger.Printf("Running command:\nsudo kernelstub -a \"%s\"", kernel_args)

	// Run the command
	_, err := command.Run("sudo", "kernelstub", "-a", kernel_args)
	errorcheck.ErrorCheck(err, "Error, kernelstub command returned exit code 1")

	// Return what we did
	return fmt.Sprintf("sudo kernelstub -a \"%s\"", kernel_args)
}

// Configures grub2 and/or systemd-boot using grubby
func Set_Grubby() string {
	// Get the config
	config := GetConfig()

	// Get the kernel args
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// Write to logger
	logger.Printf("Running command:\nsudo grubby --update-kernel=ALL --args=\"%s\"", kernel_args)

	// Run the command
	_, err := command.Run("sudo", "grubby", "--update-kernel=ALL", fmt.Sprintf("--args=%s", kernel_args))
	errorcheck.ErrorCheck(err, "Error, grubby command returned exit code 1")

	// Return what we did
	return fmt.Sprintf("sudo grubby --update-kernel=ALL --args=\"%s\"", kernel_args)
}

func Configure_Grub2() {
	// Get the config struct
	config := GetConfig()

	// Make the config file path
	conffile := fmt.Sprintf("%s/grub", config.Path.DEFAULT)

	// Make sure we start from scratch by deleting any old file
	if fileio.FileExist(conffile) {
		os.Remove(conffile)
	}

	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)
	sysfile := syspath_re.ReplaceAllString(conffile, "")

	// Make a regex to find the LINUX lines
	cmdline_default_re := regexp.MustCompile(`^GRUB_CMDLINE_LINUX_DEFAULT=\"(.+)\"$`)
	currentargs_re := regexp.MustCompile(`^GRUB_CMDLINE_LINUX(_DEFAULT|)=\"(.?|.+)\"$`)

	// Make a bool so we know if we edited the default line if both are in the template
	default_edited := false

	// Read the mkinitcpio file
	grub_content := fileio.ReadLines(sysfile)

	// Write to logger
	logger.Printf("Read %s:\n%s", sysfile, strings.Join(grub_content, "\n"))

	for _, line := range grub_content {
		if currentargs_re.MatchString(line) {
			// Get the current modules
			old_kernel_args := strings.Split(currentargs_re.ReplaceAllString(line, "${2}"), " ")

			// Clean up the old arguments by removing vfio related kernel arguments
			new_kernel_args := clean_Grub2_Args(old_kernel_args)

			// Get the kernel args from our config
			kernel_args := fileio.ReadFile(config.Path.CMDLINE)

			// Add our kernel args to the list
			new_kernel_args = append(new_kernel_args, kernel_args)

			// If we are at the line starting with MODULES=
			if cmdline_default_re.MatchString(line) {
				// Write to logger
				logger.Printf("Replacing line in %s:\n%s\nWith:\nGRUB_CMDLINE_LINUX_DEFAULT=\"%s\"\n", conffile, line, strings.Join(new_kernel_args, " "))

				// Write the modules line we generated
				fileio.AppendContent(fmt.Sprintf("GRUB_CMDLINE_LINUX_DEFAULT=\"%s\"\n", strings.Join(new_kernel_args, " ")), conffile)

				// Mark the default line as edited so we can skip the non default line
				default_edited = true
			} else {
				// If we have not edited the GRUB_CMDLINE_LINUX_DEFAULT line
				if !default_edited {
					// Write to logger
					logger.Printf("Replacing line in %s:\n%s\nWith:\nGRUB_CMDLINE_LINUX=\"%s\"\n", conffile, line, strings.Join(new_kernel_args, " "))

					// Write the modules line we generated
					fileio.AppendContent(fmt.Sprintf("GRUB_CMDLINE_LINUX=\"%s\"\n", strings.Join(new_kernel_args, " ")), conffile)
				}
			}
		} else {
			// Write the line to the file since it does not match our regex
			fileio.AppendContent(fmt.Sprintf("%s\n", line), conffile)
		}
	}
}

func clean_Grub2_Args(old_kernel_args []string) []string {
	// Make a regex to get the VFIO related kernel arguments removed, if they already existed
	vfio_args_re := regexp.MustCompile(`(amd|intel)_iommu=(on|1)|iommu=(pt|on)|vfio_pci.ids=.+|vfio_pci.disable_vga=\d{1}`)

	// Make a stringlist to keep our new arguments
	var clean_kernel_args []string

	// Loop through current kernel_args and add anything that isnt vfio or vendor-reset related
	for _, v := range old_kernel_args {
		// If what we find is not a vfio module or vendor-reset module
		if !vfio_args_re.MatchString(v) {
			// Add module to module list
			clean_kernel_args = append(clean_kernel_args, v)
		}
	}

	// Return cleaned up arguments
	return clean_kernel_args
}

// This function copies our config to /etc/default/grub and updates grub
func Set_Grub2() ([]string, error) {
	// Get the config
	config := GetConfig()

	// Get the conf file
	conffile := fmt.Sprintf("%s/grub", config.Path.DEFAULT)

	// Write to logger
	logger.Printf("Executing command:\nsudo cp -v \"%s\" /etc/default/grub", conffile)

	// Since we should be elevated with our sudo token we will copy with cp
	// (using built in functions will not work as we are running as the normal user)
	output, err := command.Run("sudo", "cp", "-v", conffile, "/etc/default/grub")
	errorcheck.ErrorCheck(err, fmt.Sprintf("Failed to copy %s to /etc/default/grub", conffile))

	// Write output to logger
	logger.Printf(strings.Join(output, "\n"))

	// Set a variable for the mkconfig command
	mkconfig := "grub-mkconfig"
	// Check for grub-mkconfig
	_, err = command.Run("which", "grub-mkconfig")
	if err == nil {
		// Set binary as grub-mkconfig
		mkconfig = "grub-mkconfig"
	} else {
		mkconfig = "grub2-mkconfig"
	}

	// Update grub.cfg
	if fileio.FileExist("/boot/grub/grub.cfg") {
		output = append(output, fmt.Sprintf("sudo %s -o /boot/grub/grub.cfg", mkconfig))
		mklog, err := command.Run("sudo", mkconfig, "-o", "/boot/grub/grub.cfg")
		logger.Printf(strings.Join(mklog, "\n"))
		errorcheck.ErrorCheck(err, "Failed to update /boot/grub/grub.cfg")
	} else {
		output = append(output, fmt.Sprintf("sudo %s -o /boot/grub/grub.cfg\nSee debug.log for more detailed output", mkconfig))
		mklog, err := command.Run("sudo", mkconfig, "-o", "/boot/grub2/grub.cfg")
		logger.Printf(strings.Join(mklog, "\n"))
		errorcheck.ErrorCheck(err, "Failed to update /boot/grub/grub.cfg")
	}

	return output, err
}
