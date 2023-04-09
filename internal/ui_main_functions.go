package internal

import (
	"regexp"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
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
		// Generate the VBIOS dumper script once the user has selected a GPU
		generateVBIOSDumper(*m)
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

		// If user selected yes then
		if selectedItem.(item).title == "YES" {
			// Add disable VFIO video to the config
			m.disableVFIOVideo()
		}

		// Configure modprobe
		configs.Set_Modprobe()

		// Go to the next view
		m.focused++

	case INTRO:
		// This is an OK Dialog
		// Create the config folder and the files related to this system
		configs.InitConfigs()

		// Go to the next view
		m.focused++

	case DONE:
		return true
	}

	return false
}

func (m *model) disableVFIOVideo() {
	// Get the config
	config := configs.GetConfig()

	// Add to the kernel arguments that we want to disable VFIO video output on the host
	fileio.AppendContent(" vfio_pci.disable_vga=1", config.Path.CMDLINE)
}
