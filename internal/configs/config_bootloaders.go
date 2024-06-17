package configs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/klauspost/cpuid/v2"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
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

// Set_KernelStub configures systemd-boot using kernelstub.
func Set_KernelStub(isRoot bool) {
	// Get the config
	config := GetConfig()

	// Get the kernel args
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// Run and log, check for errors
	common.ErrorCheck(command.ExecAndLogSudo(isRoot, true,
		"kernelstub -a "+kernel_args,
	),
		"Error, kernelstub command returned exit code 1",
	)
}

// Set_Grubby configures grub2 and/or systemd-boot using grubby
func Set_Grubby(isRoot bool) string {
	// Get the config
	config := GetConfig()

	// Get the kernel args
	kernel_args := fileio.ReadFile(config.Path.CMDLINE)

	// Run and log, check for errors
	err := command.ExecAndLogSudo(isRoot, true, "grubby --update-kernel=ALL "+fmt.Sprintf("--args=%s", kernel_args))
	common.ErrorCheck(err, "Error, grubby command returned exit code 1")

	// Return what we did
	return fmt.Sprintf("Executed: sudo grubby --update-kernel=ALL --args=\"%s\"", kernel_args)
}

func Configure_Grub2() {
	// Get the config struct
	config := GetConfig()

	// Make the config file path
	conffile := fmt.Sprintf("%s/grub", config.Path.DEFAULT)

	// Make sure we start from scratch by deleting any old file
	if exists, _ := fileio.FileExist(conffile); exists {
		_ = os.Remove(conffile)
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
	logger.Printf("Read %s:\n%s\n", sysfile, strings.Join(grub_content, "\n"))

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
				} else {
					// Since we have edited the GRUB_CMDLINE_LINUX_DEFAULT line, we will just clean up the non default line
					fileio.AppendContent(fmt.Sprintf("GRUB_CMDLINE_LINUX=\"%s\"\n", strings.Join(clean_Grub2_Args(old_kernel_args), " ")), conffile)
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
	vfio_args_re := regexp.MustCompile(`(amd|intel)_iommu=(on|1)|iommu=(pt|on)|vfio_pci.ids=.+|vfio_pci.disable_vga=\d{1}|rd.driver.pre=vfio_pci`)

	// Make a stringlist to keep our new arguments
	var clean_kernel_args []string

	// Loop through current kernel_args and add anything that isnt vfio or vendor-reset related
	for _, v := range old_kernel_args {
		// If what we find is not a vfio argument
		if !vfio_args_re.MatchString(v) {
			// Add argument to the list
			clean_kernel_args = append(clean_kernel_args, v)
		}
	}

	// Return cleaned up arguments
	return clean_kernel_args
}

// Set_Grub2 copies our config to /etc/default/grub and updates grub
func Set_Grub2(isRoot bool) error {
	// Get the config
	config := GetConfig()

	// Get the conf file
	conffile := fmt.Sprintf("%s/grub", config.Path.DEFAULT)

	// Get the sysfile
	sysfile_re := regexp.MustCompile(`^config`)
	sysfile := sysfile_re.ReplaceAllString(conffile, "")

	// [CopyToSystem] will log the operation
	// logger.Printf("Executing command:\nsudo cp -v \"%s\" %s\n", conffile, sysfile)

	// Copy files to system, logging and error checking is done in the function
	CopyToSystem(isRoot, conffile, sysfile)

	// Set a variable for the mkconfig command
	var mkconfig string
	var grubPath = "/boot/grub/grub.cfg"
	var lpErr error

	// Check for grub-mkconfig
	mkconfig, lpErr = exec.LookPath("grub-mkconfig")
	switch {
	case errors.Is(lpErr, exec.ErrNotFound) || mkconfig == "":
		// Check for grub2-mkconfig
		mkconfig, lpErr = exec.LookPath("grub2-mkconfig")
		if lpErr == nil && mkconfig != "" {
			grubPath = "/boot/grub2/grub.cfg"
			break // skip below, we found grub2-mkconfig
		}
		if lpErr == nil {
			// we know mkconfig is empty despite no error;
			// so set an error for [common.ErrorCheck].
			lpErr = errors.New("neither grub-mkconfig or grub2-mkconfig found")
		}
		common.ErrorCheck(lpErr, lpErr.Error()+"\n")
		return lpErr // note: unreachable as [common.ErrorCheck] calls fatal
	default:
	}

	_, mklog, err := command.RunErrSudo(isRoot, mkconfig, "-o", grubPath)

	// tabulate the output, [command.RunErrSudo] logged the execution.
	logger.Printf("\t" + strings.Join(mklog, "\n\t"))
	common.ErrorCheck(err, "Failed to update /boot/grub/grub.cfg")

	// always returns nil as [common.ErrorCheck] calls fatal
	// keeping the ret signature, as we should consider passing down errors
	// but that's a massive rabbit hole to go down for this codebase as a whole
	return err
}
