package pages

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gookit/color"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/HikariKnight/quickpassthrough/internal/lsiommu"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
)

func SelectGPU(config *configs.Config) {
	// Clear the screen
	command.Clear()

	// Get the users GPUs
	gpus := lsiommu.GetIOMMU("-g", "-F", "vendor:,prod_name,optional_revision:,device_id")

	// Generate a list of choices based on the GPUs and get the users selection
	choice := menu.GenIOMMUMenu("Select a GPU to view the IOMMU groups of", gpus)

	// Parse the choice
	switch choice {
	case "back":
		Welcome()
	case "":
		// If ESC is pressed
		fmt.Println("")
		os.Exit(0)
	default:
		config.Gpu_Group = choice
		viewGPU(config)
	}
}

func viewGPU(config *configs.Config, ext ...int) {
	// Clear the screen
	command.Clear()

	// Set mode to relative
	mode := "-r"

	// Set mode to relative extended
	if len(ext) > 0 {
		mode = "-rr"
	}

	// Get the IOMMU listings for GPUs
	group := lsiommu.GetIOMMU("-g", mode, "-i", config.Gpu_Group, "-F", "vendor:,prod_name,optional_revision:,device_id")

	// Write a title
	title := color.New(color.BgHiBlue, color.White, color.Bold)
	title.Println("This list should only show devices related to your GPU (usually 1 video, 1 audio device)")

	// Print all the gpus
	for _, v := range group {
		fmt.Println(v)
	}

	// Add a new line for tidyness
	fmt.Println("")

	// Make an empty string
	var choice string

	// Change choices depending on if we have done an extended search or not
	if len(ext) > 0 {
		choice = menu.YesNoManual("Use this GPU (any extra devices listed may or may not be linked to it) for passthrough?")
	} else {
		choice = menu.YesNoEXT("Use this GPU (and related devices) for passthrough?")
	}

	// Parse the choice
	switch choice {
	case "":
		// If ESC is pressed
		fmt.Println("")
		os.Exit(0)

	case "ext":
		// Run an extended relative search
		viewGPU(config, 1)

	case "n":
		// Go back to selecting a gpu
		SelectGPU(config)

	case "y":
		// Get the device ids for the selected gpu using ls-iommu
		config.Gpu_IDs = lsiommu.GetIOMMU("-g", mode, "-i", config.Gpu_Group, "--id")

	case "manual":
		config.Gpu_IDs = menu.ManualInput(
			"Please manually enter the vendorID:deviceID for every device to use except PCI Express Switches\n"+
				"NOTE: All devices sharing the same IOMMU group will still get pulled into the VM!",
			"xxxx:yyyy,xxxx:yyyy,xxxx:yyyy",
		)
	}

	logger.Printf("Checking for duplicate device Ids")
	hasDuplicateDeviceIds := detectDuplicateDeviceIds(config.Gpu_Group, config.Gpu_IDs)

	if hasDuplicateDeviceIds {
		config.HasDuplicateDeviceIds = true
		config.Gpu_Addresses = lsiommu.GetIOMMU("-g", mode, "-i", config.Gpu_Group, "--pciaddr")
	}

	// If the kernel_args file already exists
	if exists, _ := fileio.FileExist(config.Path.CMDLINE); exists {
		// Delete it as we will have to make a new one anyway
		err := os.Remove(config.Path.CMDLINE)
		common.ErrorCheck(err, fmt.Sprintf("Could not remove %s", config.Path.CMDLINE))
	}

	// Write initial kernel_arg file
	configs.Set_Cmdline(config.Gpu_IDs, !config.HasDuplicateDeviceIds)

	// Go to the vbios dumper page
	genVBIOS_dumper(config)
}

func detectDuplicateDeviceIds(selectedGpuGroup string, selectedDeviceIds []string) bool {
	// TODO: this would be made much simpler if ls-iommu allowed using the --id flag without
	// the "-i" flag.
	gpus := lsiommu.GetIOMMU("-g", "-F", "vendor:,prod_name,optional_revision:,device_id")
	iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
	iommuGroups := []string{}
	for _, gpu := range gpus {
		iommuGroup := iommu_group_regex.FindString(gpu)
		iommuGroups = append(iommuGroups, iommuGroup)
	}

	allDeviceIds := []string{}
	for _, group := range iommuGroups {
		if group == selectedGpuGroup {
			continue
		}

		deviceIds := lsiommu.GetIOMMU("-g", "-r", "-i", group, "--id")
		for _, deviceId := range deviceIds {
			allDeviceIds = append(allDeviceIds, deviceId)
		}
	}

	for _, deviceId := range allDeviceIds {
		for _, selectedDeviceId := range selectedDeviceIds {
			if deviceId == selectedDeviceId {
				logger.Printf("Found duplicate device id: %s", deviceId)
				return true
			}
		}
	}

	return false
}
