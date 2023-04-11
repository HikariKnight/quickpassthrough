package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/HikariKnight/quickpassthrough/pkg/uname"
)

// This function processes the enter event
func (m *model) processSelection() bool {
	switch m.focused {
	case GPUS:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Add the gpu group to our model (this is so we can grab the vbios details later)
		m.gpu_group = iommu_group

		// Get all the gpu devices and related devices (same device id or in the same group)
		items := iommuList2ListItem(getIOMMU("-grr", "-i", m.gpu_group, "-F", "vendor:,prod_name,optional_revision:,device_id"))

		// Add the devices to the list
		m.lists[GPU_GROUP].SetItems(items)

		// Change focus to next index
		m.focused++

	case GPU_GROUP:
		// Get the config
		config := configs.GetConfig()

		// Get the vbios path
		m.vbios_path = getIOMMU("-g", "-i", m.gpu_group, "--rom")[0]

		// Generate the VBIOS dumper script once the user has selected a GPU
		configs.GenerateVBIOSDumper(m.vbios_path)

		// Get the device ids for the selected gpu using ls-iommu
		m.gpu_IDs = getIOMMU("-gr", "-i", m.gpu_group, "--id")

		// If the kernel_args file already exists
		if fileio.FileExist(config.Path.CMDLINE) {
			// Delete it as we will have to make a new one anyway
			err := os.Remove(config.Path.CMDLINE)
			errorcheck.ErrorCheck(err, fmt.Sprintf("Could not remove %s", config.Path.CMDLINE))
		}

		// Write initial kernel_arg file
		configs.Set_Cmdline(m.gpu_IDs)

		// Change focus to the next view
		m.focused++

	case USB:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Get the USB controllers in the selected iommu group
		items := iommuList2ListItem(getIOMMU("-ur", "-i", iommu_group, "-F", "vendor:,prod_name,optional_revision:,device_id"))

		// Add the items to the list
		m.lists[USB_GROUP].SetItems(items)

		// Change focus to next index
		m.focused++

	case USB_GROUP:
		m.focused++

	case VBIOS:
		// This is just an OK Dialog
		m.focused++

	case VIDEO:
		// This is a YESNO Dialog
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Get our config struct
		config := configs.GetConfig()

		// If user selected yes then
		if selectedItem.(item).title == "YES" {
			// Add disable VFIO video to the config
			configs.DisableVFIOVideo(1)
		} else {
			// Add disable VFIO video to the config
			configs.DisableVFIOVideo(0)
		}

		// If we have files for modprobe
		if fileio.FileExist(config.Path.MODPROBE) {
			// Configure modprobe
			configs.Set_Modprobe(m.gpu_IDs)
		}

		// If we have a folder for dracut
		if fileio.FileExist(config.Path.DRACUT) {
			// Configure dracut
			configs.Set_Dracut()
		}

		// If we have a mkinitcpio.conf file
		if fileio.FileExist(config.Path.MKINITCPIO) {
			configs.Set_Mkinitcpio()
		}

		// Configure grub2 here as we can make the config without sudo
		if config.Bootloader == "grub2" {
			// Write to logger
			logger.Printf("Configuring grub2 manually")
			configs.Configure_Grub2()
		}

		// Go to the next view
		//m.focused++

		// Because we have no QuickEmu support yet, just skip USB Controller configuration
		m.focused = INSTALL
		return true

	case INTRO:
		// This is an OK Dialog
		// Create the config folder and the files related to this system
		configs.InitConfigs()

		// Go to the next view
		m.focused++

	case DONE:
		// Return true so that the application will exit nicely
		return true
	}

	// Return false as we are not done
	return false
}

// This function starts the install process
// It takes 1 auth string as variable
func (m *model) install() {
	// Get the config
	config := configs.GetConfig()

	// Make a stringlist to keep the output to show the user
	var output []string

	// Based on the bootloader, setup the configuration
	if config.Bootloader == "kernelstub" {
		// Write to logger
		logger.Printf("Configuring systemd-boot using kernelstub")

		// Configure kernelstub
		output = append(output, configs.Set_KernelStub())

	} else if config.Bootloader == "grubby" {
		// Write to logger
		logger.Printf("Configuring bootloader using grubby")

		// Configure kernelstub
		output = append(output, configs.Set_Grubby())

	} else if config.Bootloader == "grub2" {
		// Write to logger
		logger.Printf("Configuring grub2 manually")
		grub_output, _ := configs.Set_Grub2()
		output = append(output, grub_output...)

	} else {
		kernel_args := fileio.ReadFile(config.Path.CMDLINE)
		logger.Printf("Unsupported bootloader, please add the below line to your bootloaders kernel arguments\n%s", kernel_args)
	}

	// A lot of linux systems support modprobe along with their own module system
	// So copy the modprobe files if we have them
	modprobeFile := fmt.Sprintf("%s/vfio.conf", config.Path.MODPROBE)
	if fileio.FileExist(modprobeFile) {
		// Copy initramfs-tools module to system
		output = append(output, configs.CopyToSystem(modprobeFile, "/etc/modprobe.d/vfio.conf"))
	}

	// Copy the config files for the system we have
	initramfsFile := fmt.Sprintf("%s/modules", config.Path.INITRAMFS)
	dracutFile := fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT)
	if fileio.FileExist(initramfsFile) {
		// Copy initramfs-tools module to system
		output = append(output, configs.CopyToSystem(initramfsFile, "/etc/initramfs-tools/modules"))

		// Copy the modules file to /etc/modules
		output = append(output, configs.CopyToSystem(config.Path.ETCMODULES, "/etc/modules"))

		// Write to logger
		logger.Printf("Executing: sudo update-initramfs -u")

		// Update initramfs
		output = append(output, "Executed: sudo update-initramfs -u\nSee debug.log for detailed output")
		cmd_out, cmd_err, _ := command.RunErr("sudo", "update-initramfs", "-u")

		cmd_out = append(cmd_out, cmd_err...)

		// Write to logger
		logger.Printf(strings.Join(cmd_out, "\n"))
	} else if fileio.FileExist(dracutFile) {
		// Copy dracut config to /etc/dracut.conf.d/vfio
		output = append(output, configs.CopyToSystem(dracutFile, "/etc/dracut.conf.d/vfio"))

		// Get systeminfo
		sysinfo := uname.New()

		// Write to logger
		logger.Printf("Executing: sudo dracut -f -v --kver %s", sysinfo.Release)

		// Update initramfs
		output = append(output, fmt.Sprintf("Executed: sudo dracut -f -v --kver %s\nSee debug.log for detailed output", sysinfo.Release))
		cmd_out, cmd_err, _ := command.RunErr("sudo", "dracut", "-f", "-v", "--kver", sysinfo.Release)

		cmd_out = append(cmd_out, cmd_err...)

		// Write to logger
		logger.Printf(strings.Join(cmd_out, "\n"))
	} else if fileio.FileExist(config.Path.MKINITCPIO) {
		// Copy dracut config to /etc/dracut.conf.d/vfio
		output = append(output, configs.CopyToSystem(config.Path.MKINITCPIO, "/etc/mkinitcpio.conf"))

		// Write to logger
		logger.Printf("Executing: sudo mkinitcpio -P")

		// Update initramfs
		output = append(output, "Executed: sudo mkinitcpio -P\nSee debug.log for detailed output")
		cmd_out, cmd_err, _ := command.RunErr("sudo", "mkinitcpio", "-P")

		cmd_out = append(cmd_out, cmd_err...)

		// Write to logger
		logger.Printf(strings.Join(cmd_out, "\n"))
	}

	m.installOutput = output
	m.focused++

}
